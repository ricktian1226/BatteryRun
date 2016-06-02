// xypanic
package xypanic

import (
    "bytes"
    "fmt"
    "guanghuan.com/xiaoyao/common/log"
    "guanghuan.com/xiaoyao/common/mail"
    "io/ioutil"
    "os"
    "runtime"
    "time"
)

//panic开关
var Panic_Switch = true

//from martini source code.
var (
    dunno     = []byte("???")
    centerDot = []byte(".")
    dot       = []byte(".")
    slash     = []byte("/")
)

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
    n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
    if n < 0 || n >= len(lines) {
        return dunno
    }
    return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
    fn := runtime.FuncForPC(pc)
    if fn == nil {
        return dunno
    }
    name := []byte(fn.Name())
    // The name includes the path name to the package, which is unnecessary
    // since the file name is already included.  Plus, it has center dots.
    // That is, we see
    //	runtime/debug.*T.ptrmethod
    // and want
    //	*T.ptrmethod
    // Also the package path might contains dot (e.g. code.google.com/...),
    // so first eliminate the path prefix
    if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
        name = name[lastslash+1:]
    }
    if period := bytes.Index(name, dot); period >= 0 {
        name = name[period+1:]
    }
    name = bytes.Replace(name, centerDot, dot, -1)
    return name
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
    buf := new(bytes.Buffer) // the returned data
    // As we loop, we open files and read them. These variables record the currently
    // loaded file.
    var lines [][]byte
    var lastFile string
    for i := skip; ; i++ { // Skip the expected number of frames
        pc, file, line, ok := runtime.Caller(i)
        if !ok {
            break
        }
        // Print this much at least.  If we can't find the source, it won't show.
        //fmt.Printf("%s:%d (0x%x)\n", file, line, pc)
        fmt.Fprintf(buf, "\t%s:%d (0x%x)\n", file, line, pc)
        if file != lastFile {
            data, err := ioutil.ReadFile(file)
            if err != nil {
                continue
            }
            lines = bytes.Split(data, []byte{'\n'})
            lastFile = file
        }
        fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
    }
    return buf.Bytes()
}

func Crash() {

    //no need to panic
    //fmt.Printf("Panic_Switch : %t", Panic_Switch)
    if !Panic_Switch {
        return
    }

    if err := recover(); err != nil {
        file := fmt.Sprintf("%s/Crash.%d.log", xylog.DefConfig.Path, xylog.DefConfig.LogId)
        fd, errCreate := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
        if errCreate != nil {
            fmt.Println("create " + file + "failed")
            return
        }
        defer fd.Close()

        fmt.Fprintf(fd, "%s panic : %v\n", time.Now().String(), err)

        stacks := string(stack(3))

        fmt.Fprintf(fd, "%v\n", stacks)

        //发送告警邮件
        subject := fmt.Sprintf("[BusinessPanic!!!] battery_%s_server_%d_%d", xylog.DefConfig.AppName, xylog.DefConfig.DCId, xylog.DefConfig.NodeId)
        message := string("===== This mail is sent by system, please don't reply =====\n") +
            stacks +
            string("\n==============\n")

        xymail.Send(subject, message, false)

        runtime.Goexit()
    }
}
