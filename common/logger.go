package common

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"hdyx/global"
	"hdyx/net/ioBuf"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

var DEFAULT_TIME_TYPE = zapcore.TimeEncoderOfLayout(time.DateTime)

type tLogger struct {
	logger *zap.Logger
}

var Logger tLogger

func (this tLogger) InfoLog(message string) {
	this.logger.Info(message)
}

func (this tLogger) SystemErrorLog(message ...any) {
	this.logger.Panic(fmt.Sprintln(message...) + string(debug.Stack()))
	runtime.Goexit()
}

func (this tLogger) GameErrorLog(conCtx *global.ConContext, errId uint32, message ...any) {
	this.logger.Error(fmt.Sprintln(message...) + string(debug.Stack()))

	out := ioBuf.OutPutBuf{
		CmdCode:        1,
		ProtocolSwitch: 0,
		CmdMerge:       conCtx.GetConGlobalVal().Cmd,
		ResponseStatus: errId,
		Data:           nil,
	}

	OutPutStream[*ioBuf.OutPutBuf](conCtx, &out)
	runtime.Goexit()
}

func init() {
	logCfg := GetYamlMapCfg("serverCfg", "logger").(map[string]any)
	maxSizeCfg := logCfg["maxSize"].(int)
	maxBackupsCfg := logCfg["maxBackups"].(int)
	maxKeepDayCfg := logCfg["maxKeepDay"].(int)
	isCompressCfg := logCfg["isCompress"].(bool)

	var coreArr []zapcore.Core

	//获取编码器
	encoderConfig := zap.NewProductionEncoderConfig() //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = DEFAULT_TIME_TYPE      //指定时间格式
	//encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	//日志级别判定
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	//info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   LOG_PATH + "/info_" + time.Now().Format(time.DateOnly) + ".log", //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    maxSizeCfg,                                                      //文件大小限制,单位MB
		MaxBackups: maxBackupsCfg,                                                   //最大保留日志文件数量
		MaxAge:     maxKeepDayCfg,                                                   //日志文件保留天数
		Compress:   isCompressCfg,                                                   //是否压缩处理
	})
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	//error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   LOG_PATH + "/error_" + time.Now().Format(time.DateOnly) + ".log", //日志文件存放目录
		MaxSize:    maxSizeCfg,                                                       //文件大小限制,单位MB
		MaxBackups: maxBackupsCfg,                                                    //最大保留日志文件数量
		MaxAge:     maxKeepDayCfg,                                                    //日志文件保留天数
		Compress:   isCompressCfg,                                                    //是否压缩处理
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	//panic文件writeSyncer
	panicFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   LOG_PATH + "/sysErr_" + time.Now().Format(time.DateOnly) + ".log", //日志文件存放目录
		MaxSize:    maxSizeCfg,                                                        //文件大小限制,单位MB
		MaxBackups: maxBackupsCfg,                                                     //最大保留日志文件数量
		MaxAge:     maxKeepDayCfg,                                                     //日志文件保留天数
		Compress:   isCompressCfg,                                                     //是否压缩处理
	})
	panicFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(panicFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
	coreArr = append(coreArr, panicFileCore)
	Logger.logger = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) //zap.AddCaller()为显示文件名和行号，可省略

}
