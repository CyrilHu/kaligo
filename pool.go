package kaligo

import (
    "time"
    //"fmt"
    //"reflect"
)

var nowFunc = time.Now // for testing

// ConnPool is ...
type ConnPool struct {
	// Dial is an application supplied function for creating and configuring a
	// connection
	Dial func() (interface{}, error)
    //Dial func() (*autorc.Conn, error)

    // Maximum number of idle connections in the pool.
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	closed bool
	active int

    idle chan interface{}
}

type idleConn struct {
	c interface{}
	t time.Time
}

// InitPool is 批量生成连接，并把连接放到连接池channel里面
func (connPoll *ConnPool) InitPool() error{
    connPoll.idle = make(chan interface{}, connPoll.MaxActive)
    for x := 0; x < connPoll.MaxActive; x++ {
        //conn, err := this.Dial()
        // 这里返回DB类，而不是返回mysql.Conn，否则DB类的insert，update 这些Active Record类方法没法使用
        db, err := connPoll.Dial()
        //fmt.Println(" --- reflect --- ", reflect.TypeOf(db))
        if err != nil {
            return err
        }
        //connPoll.idle <-conn
        connPoll.idle <-idleConn{t: nowFunc(), c: db}
    }
    return nil
}

// Get is 从连接池里取出连接
func (connPoll *ConnPool) Get() interface{} {
    // 如果空闲连接为空，初始化连接池
    if connPoll.idle == nil {
		err := connPoll.InitPool()
		if err != nil {
			return nil
		}
    }

    // 赋值一下好给下面回收和返回
    //conn := <-connPoll.idle
    //idleConn
    ic := <-connPoll.idle
    // 这里要用 (idleConn) 把interface{} 类型转化为 idleConn 类型的，否则拿不到里面的属性t、c
    conn := ic.(idleConn).c
    //fmt.Println(conn.(*DB).conn)
    //fmt.Println(" --- reflect --- ", reflect.TypeOf(conn))

    // 使用完把连接回收到连接池里
    defer connPoll.Release(conn)

    // 因为channel是有锁的，所以就没必要借助sync.Mutex来进行读写锁定
    // container/list就需要锁住，不然并发就互抢出问题了
    return conn
}

// Release is 回收连接到连接池
func (connPoll *ConnPool) Release(conn interface{}) {
    //connPoll.idle <-conn
    connPoll.idle <-idleConn{t: nowFunc(), c: conn}
}

