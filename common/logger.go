package common

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"reflect"
	"runtime"
	"time"
	"unsafe"
)

var DEFAULT_TIME_TYPE = zapcore.TimeEncoderOfLayout(time.DateTime)

type tLogger struct {
	logger *zap.Logger
}

var Logger tLogger

func (this tLogger) InfoLog(message string) {
	this.logger.Info(message)
}

func (this tLogger) ErrorLog(message string) {
	this.logger.Error(message)
	fmt.Println(callStack())
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

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
	Logger.logger = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) //zap.AddCaller()为显示文件名和行号，可省略
}

func callStack() string {
	var trace string
	pcs := make([]uintptr, 32)
	n := runtime.Callers(0, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		fn := runtime.FuncForPC(frame.PC)
		file, line := fn.FileLine(frame.PC)
		args := getArgs(frame)
		fmt.Printf("%s - %s(%v)\n\t%s:%d\n", trace, fn.Name(), args, file, line)
		if !more {
			break
		}
	}
	return trace
}

func getArgs(frame runtime.Frame) []interface{} {
	fn := frame.Func
	f := reflect.ValueOf(fn)
	if f.Kind() != reflect.Func { // 检查是否为函数类型
		return nil
	}
	in := make([]reflect.Value, f.Type().NumIn())
	for i := range in {
		// 获取参数类型
		t := f.Type().In(i)

		// 如果参数类型是 interface{}，则返回 nil
		if t.Kind() == reflect.Interface {
			return nil
		}

		// 创建一个具有指定类型并初始化为零值的新变量
		v := reflect.New(t).Elem()

		// 将参数值存储到新变量中
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intSize := int(unsafe.Sizeof(int(0)))
			ptr := unsafe.Pointer(&frame.Entry)
			val := *(*int64)(unsafe.Pointer(uintptr(ptr) + uintptr(i*intSize)))
			v.SetInt(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			intSize := int(unsafe.Sizeof(int(0)))
			ptr := unsafe.Pointer(&frame.Entry)
			val := *(*uint64)(unsafe.Pointer(uintptr(ptr) + uintptr(i*intSize)))
			v.SetUint(val)
		case reflect.Float32, reflect.Float64:
			floatSize := int(unsafe.Sizeof(float32(0)))
			ptr := unsafe.Pointer(&frame.Entry)
			val := *(*float64)(unsafe.Pointer(uintptr(ptr) + uintptr(i*floatSize)))
			v.SetFloat(val)
		case reflect.String:
			intSize := int(unsafe.Sizeof(int(0)))
			ptr := unsafe.Pointer(&frame.Entry)
			strPtr := *(*uintptr)(unsafe.Pointer(uintptr(ptr) + uintptr(i*intSize))) // 获取字符串指针
			str := *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: strPtr}))   // 使用 StringHeader 解析字符串
			v.SetString(str)
		default:
			// 如果参数类型是自定义类型，则无法获取其值
			return nil
		}

		in[i] = v
	}
	return getInterfaceSlice(in)
}

func getInterfaceSlice(values []reflect.Value) []interface{} {
	out := make([]interface{}, len(values))
	for i, v := range values {
		out[i] = v.Interface()
	}
	return out
}
