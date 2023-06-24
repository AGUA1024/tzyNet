package tServer

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
	"tzyNet/tCommon"
)

type DbOperator struct {
	pool     dbPool
	database string
}

type dbPool struct {
	db           *gorm.DB           // GORM库的DB对象
	maxOpenConns int                // 最大打开连接数
	maxIdleConns int                // 最大空闲连接数
	connLifetime time.Duration      // 连接存活时间
	freeConn     []*gorm.DB         // 空闲连接列表
	mutex        sync.Mutex         // 互斥锁，用于保护共享变量
	connRequests chan chan *gorm.DB // 连接请求通道，用于处理新连接的生成和获取
}

type dbCfgObj struct {
	User         string
	Pass         string
	AllHostCfg   map[string]any
	MaxOpenCon   int
	MaxIdleConns int
	ConLiveTime  int
	PieceNum     int
}

var (
	dbCfg      dbCfgObj
	arrDbPools []*dbPool
)

var mpDatabase = map[string]string{
	"game": "hdyx_game",
}

func init() {
	//// 配置初始化
	//dbCfgIni()
	//// 数据库初始化
	//dbPoolInit()
	fmt.Println("--Mysql初始化完成")
}

func GetDb(uid uint64, database string) DbOperator {
	dbName, ok := mpDatabase[database]
	if !ok {
		tCommon.Logger.SystemErrorLog("DB_NAME_NOT_FOUNT:" + database)
	}

	piece := uid % uint64(dbCfg.PieceNum)

	dataBase := fmt.Sprintf("%s%02d", dbName+"_", piece+1)
	var dbOp = DbOperator{
		pool:     *arrDbPools[piece],
		database: dataBase,
	}

	return dbOp
}

// 查询数据
func (op *DbOperator) QueryData(ctx context.Context, table string, where map[string]interface{}, dest interface{}) error {
	db, err := op.pool.getDbFromPool()
	if err != nil {
		return err
	}
	defer op.pool.releaseConn(db)

	dbTable := op.database + "." + table
	result := db.WithContext(ctx).Table(dbTable).Where(where).Find(dest)

	return result.Error
}

// 插入数据
func (op *DbOperator) InsertData(ctx context.Context, table string, data interface{}) error {
	db, err := op.pool.getDbFromPool()
	if err != nil {
		return err
	}
	defer op.pool.releaseConn(db)

	dbTable := op.database + "." + table
	result := db.WithContext(ctx).Table(dbTable).Create(data)
	return result.Error
}

// 更新数据
func (op *DbOperator) UpdateData(ctx context.Context, table string, where map[string]interface{}, data interface{}) error {
	db, err := op.pool.getDbFromPool()
	if err != nil {
		return err
	}
	defer op.pool.releaseConn(db)

	dbTable := op.database + "." + table
	result := db.WithContext(ctx).Table(dbTable).Where(where).Updates(data)
	return result.Error
}

// 删除数据
func (op *DbOperator) DeleteData(ctx context.Context, table string, where map[string]interface{}) error {
	db, err := op.pool.getDbFromPool()
	if err != nil {
		return err
	}
	defer op.pool.releaseConn(db)

	dbTable := op.database + "." + table
	result := db.WithContext(ctx).Table(dbTable).Where(where).Delete(nil)
	return result.Error
}

func dbCfgIni() {
	dbCfg.User = tCommon.GetYamlMapCfg("mysqlCfg", "mysql", "common", "user").(string)
	dbCfg.Pass = tCommon.GetYamlMapCfg("mysqlCfg", "mysql", "common", "pass").(string)
	dbCfg.AllHostCfg = tCommon.GetYamlMapCfg("mysqlCfg", "mysql", "game").(map[string]any)
	dbCfg.MaxOpenCon = tCommon.GetYamlMapCfg("mysqlCfg", "mysql", "common", "maxOpenCon").(int)
	dbCfg.MaxIdleConns = tCommon.GetYamlMapCfg("mysqlCfg", "mysql", "common", "maxIdleConns").(int)
	dbCfg.ConLiveTime = tCommon.GetYamlMapCfg("mysqlCfg", "mysql", "common", "conLiveTime").(int)
	dbCfg.PieceNum = tCommon.GetYamlMapCfg("mysqlCfg", "mysql", "common", "pieceNum").(int)
}

func dbPoolInit() {
	for i := 1; i <= dbCfg.PieceNum; i++ {
		host := fmt.Sprintf("%s%d", "host", i)
		var hostCfg map[string]any
		hostCfg = dbCfg.AllHostCfg[host].(map[string]any)

		ip := hostCfg["ip"].(string)
		port := hostCfg["port"]

		dbUrl := fmt.Sprintf("%s%d%s", dbCfg.User+":"+dbCfg.Pass+"@tcp("+ip+":", port, ")/?charset=utf8mb4")

		pool, err := newDbPool(dbUrl, dbCfg.MaxOpenCon, dbCfg.MaxIdleConns, time.Duration(dbCfg.ConLiveTime)*time.Second)
		if err != nil {
			tCommon.Logger.SystemErrorLog(fmt.Sprintln(err))
		}
		arrDbPools = append(arrDbPools, pool)
	}
}

// 创建新的连接池对象
func newDbPool(dsn string, maxOpenConns, maxIdleConns int, connLifetime time.Duration) (*dbPool, error) {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	dbSQL, err := db.DB()
	if err != nil {
		fmt.Println(3)
		return nil, err
	}
	dbSQL.SetMaxOpenConns(maxOpenConns)    // 设置最大打开连接数
	dbSQL.SetMaxIdleConns(maxIdleConns)    // 设置最大空闲连接数
	dbSQL.SetConnMaxLifetime(connLifetime) // 设置连接存活时间

	// 初始化连接池对象
	pool := &dbPool{
		db:           db,
		maxOpenConns: maxOpenConns,
		maxIdleConns: maxIdleConns,
		connLifetime: connLifetime,
		freeConn:     make([]*gorm.DB, 0),
		connRequests: make(chan chan *gorm.DB),
	}

	// 启动维护空闲连接的goroutine
	go pool.maintainFreeConn(dsn)

	return pool, nil
}

// 获取连接
func (pool *dbPool) getDbFromPool() (*gorm.DB, error) {
	// 先从空闲连接列表中获取一个连接，如果没有则尝试新建连接
	pool.mutex.Lock()
	if len(pool.freeConn) > 0 { // 如果有可用连接
		conn := pool.freeConn[0]
		pool.freeConn = pool.freeConn[1:] // 从空闲连接列表中移除该连接
		pool.mutex.Unlock()
		return conn, nil
	} else { // 如果没有可用的连接
		req := make(chan *gorm.DB)
		pool.connRequests <- req // 将请求加入连接请求通道
		pool.mutex.Unlock()
		select {
		case conn := <-req:
			return conn, nil
		case <-time.After(pool.connLifetime):
			return nil, fmt.Errorf("connection timeout")
		}
	}
}

// 释放连接
func (pool *dbPool) releaseConn(conn *gorm.DB) {
	// 如果空闲连接数量已达到最大值，则关闭该连接；
	// 否则将该连接放入空闲连接列表中
	pool.mutex.Lock()
	if len(pool.freeConn) >= pool.maxIdleConns { // 如果空闲连接数量已达到最大值
		pool.mutex.Unlock()
		db, _ := conn.DB()
		db.Close() // 关闭该连接
	} else { // 否则将该连接放入空闲连接列表中
		pool.freeConn = append(pool.freeConn, conn)
		pool.mutex.Unlock()
	}
}

// 维护空闲连接列表
func (pool *dbPool) maintainFreeConn(dsn string) {
	for {
		select {
		case req := <-pool.connRequests: // 有新的连接请求
			var conn *gorm.DB
			var err error
			if len(pool.freeConn) > 0 { // 如果有可用连接
				conn = pool.freeConn[0]
				pool.freeConn = pool.freeConn[1:] // 从空闲连接列表中移除该连接
			} else { // 如果没有可用的连接
				conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) // 新建连接
				if err != nil {                                        // 如果连接失败
					req <- nil // 返回nil到请求通道
					continue
				}
				dbSQL, _ := conn.DB()
				dbSQL.SetMaxOpenConns(pool.maxOpenConns)
				dbSQL.SetMaxIdleConns(pool.maxIdleConns)
				dbSQL.SetConnMaxLifetime(pool.connLifetime)
			}
			req <- conn
		}
	}
}
