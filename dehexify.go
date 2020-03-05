package main

import (
    "encoding/hex"
    "log"
    "regexp"
    "strings"
    
    "github.com/lxn/walk"
    . "github.com/lxn/walk/declarative"
)

func RegexReplaceWithACallback(re *regexp.Regexp, str string, callback func([]string) string) string {
    result := ""
    lastIndex := 0
    for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
        groups := []string{}
        for i := 0; i < len(v); i += 2 {
            groups = append(groups, str[v[i]:v[i+1]])
        }
        result += str[lastIndex:v[0]] + callback(groups)
        lastIndex = v[1]
    }
    return result + str[lastIndex:]
}

func main() {
    var inTE, outTE *walk.TextEdit

    MainWindow{
        Title:   "dehexify",
        Size:    Size{800, 500},
        Layout:  VBox{},
        Children: []Widget{
            HSplitter{
                Children: []Widget{
                    TextEdit{AssignTo: &inTE},
                    TextEdit{AssignTo: &outTE, ReadOnly: true},
                },
            },
            PushButton{
                Text: "de-hexify",
                OnClicked: func() {
                    re := regexp.MustCompile("%([0-9A-F]{2})") // %3C, %FD, ...
                    result := RegexReplaceWithACallback(re, inTE.Text(), func(groups []string) string {
                        decoded, _ := hex.DecodeString(groups[1])
                        return string(decoded)
                    })
                    
                    // "+" usually stands for a space.
                    outTE.SetText(strings.ReplaceAll(result, "+", " "))
                },
            },
            PushButton{
                Text: "copy output",
                OnClicked: func() {
                    if err := walk.Clipboard().SetText(outTE.Text()); err != nil {
                        log.Fatal("Copy error: ", err)
                    }
                },
            },
        },
    }.Run()
}