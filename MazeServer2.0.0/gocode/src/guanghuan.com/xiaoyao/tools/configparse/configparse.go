// business
package configParse

import (
    "bufio"

    xylog "guanghuan.com/xiaoyao/common/log"
    "os"
    "strings"

    //"github.com/astaxie/beego"
)

const TimeFormat string = "2006/1/2 15:04"

type ConfigRow (map[string]string)

func GetConfigRowListFromFileName(fileName string, configData *[]ConfigRow) (err error) {

    file, err := os.Open(fileName)
    defer file.Close()

    if err != nil {
        xylog.Error("Open file %s failed", fileName)
        return
    }

    var columNameList []string

    rb := bufio.NewReader(file)
    var lineIndex int = 0
    var lastParaNum int
    var paraNum int
    for ; ; lineIndex++ {
        baseline, err := rb.ReadString('\n')
        xylog.Debug("baseline : %v", baseline)
        if err != nil {
            return err
        }

        if len(baseline) <= 0 { //跳过无效行 和 空行
            continue
        }
        var line string = baseline
        xylog.DebugNoId("baseline 0 : %d", baseline[0])
        if lineIndex == 0 && baseline[0] == 239 { //utf8编码（去掉最开始的文件头编码“EF”）
            //line = beego.Substr(baseline, 1, len(baseline)-3) //去掉文件头
            line = Substr(baseline, 1, len(baseline)-3) //去掉文件头
            xylog.DebugNoId("line = baseline[3:]")
        }
        xylog.DebugNoId("line 0 : %d", line[0])
        xylog.DebugNoId("line : %v", line)
        subs := strings.Split(line, "\t")
        lastParaNum = paraNum
        if lastParaNum != 0 && lastParaNum != paraNum {
            return err
        }

        paraNum := len(subs)
        if paraNum <= 0 {
            xylog.ErrorNoId("subs  paraNum == %d", paraNum)
            continue
        }

        xylog.DebugNoId("subs : %d", paraNum)
        xylog.DebugNoId("subs : %v", subs)

        if len(line) > 0 && lineIndex == 0 {
            //读取字段名
            for _, sub := range subs {
                columNameList = append(columNameList, strings.TrimSpace(sub))
                xylog.DebugNoId("sub : [%s]", strings.TrimSpace(sub))
            }

        } else if len(columNameList) > 0 {
            //读取字段
            var configRow ConfigRow
            configRow = make(ConfigRow)
            for index, sub := range subs {
                configRow[columNameList[index]] = strings.TrimSpace(sub)
            }
            *configData = append(*configData, configRow)
        } else {
            xylog.ErrorNoId("columNameList  len <= 0")
        }
    }
    return err
}

// Substr returns the substr from start to length.
func Substr(s string, start, length int) string {
    bt := []rune(s)
    if start < 0 {
        start = 0
    }
    var end int
    if (start + length) > (len(bt) - 1) {
        end = len(bt)
    } else {
        end = start + length
    }
    return string(bt[start:end])
}
