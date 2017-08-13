package main

import(
    "fmt"
    "strings"
    "os"
    "golang.org/x/net/html"
    "net/http"
)

type TypeMember struct {
    Field       string
    Type        string
    Description string
}

func (tm TypeMember) String() string {
    return fmt.Sprintf("Field = %s, Type = %s, Description = %s", tm.Field, tm.Type, tm.Description)
}

func main() {
    if len(os.Args) != 2 {
        fmt.Printf("usage: %s [type]\n", os.Args[0])
        os.Exit(1)
    }

    // the url for the documentation
    url := "https://core.telegram.org/bots/api"

    // retrieve the documentation in html format
    resp, _ := http.Get(url)

    // make a tokenizer
    tokenizer := html.NewTokenizer(resp.Body)

    // all members of the desired type
    var tms []TypeMember

    // search for an <a> tag with href starting with a #
    done := false
    for !done {
        /* get the token type
         * https://godoc.org/golang.org/x/net/html#TokenType
         */
        tt := tokenizer.Next()
        switch {
        case tt == html.ErrorToken:
            // end of the document
            done = true
        case tt == html.StartTagToken:
            // check if start token is for an anchor
            tok := tokenizer.Token()
            if tok.Data == "a" {
                // check if there is name attribute with the desired value
                for _, attr := range tok.Attr {
                    if attr.Key == "name" && attr.Val == strings.ToLower(os.Args[1]) {
                        // parse next
                        tms = parseTable(tokenizer)
                    }
                }
            }
        }
    }

    formatTypeMembers(tms);
}

func parseTable(tokenizer *html.Tokenizer) []TypeMember {
    tm := make([]TypeMember, 0)

    // keep track of the first row which is a table header
    firstRowFound := false

    // search for <tr> elements and one </table>
    for {
        tt := tokenizer.Next()
        switch {
        case tt == html.StartTagToken:
            tok := tokenizer.Token()
            if tok.Data == "tr" {
                // we found a table row, check if it's first
                if firstRowFound == false {
                    firstRowFound = true
                } else {
                    tm = append(tm, parseRow(tokenizer))
                }
            }
        case tt == html.EndTagToken:
            tok := tokenizer.Token()
            if tok.Data == "table" {
                // we are finished
                return tm
            }
        }
    }
}

func parseRow(tokenizer *html.Tokenizer) TypeMember {
    // store the type member data
    tm := TypeMember{}

    // keep track of the member
    tdidx := 0

    // search for 3 <td> and one </tr>
    done := false
    for !done {
        tt := tokenizer.Next()
        tok := tokenizer.Token()

        if tt == html.TextToken && tok.Data != "\n" {
            switch tdidx {
            case 0:
                tm.Field += tok.Data
            case 1:
                tm.Type += tok.Data
            case 2:
                tm.Description += tok.Data
            }
        }

        if tt == html.EndTagToken {
            switch tok.Data {
            case "td":
                tdidx++
            case "tr":
                done = true
            }
        }
    }

    return tm
}

func formatTypeMembers(tms []TypeMember) {
    fmt.Printf("// https://core.telegram.org/bots/api#%s\n", os.Args[1])
    fmt.Printf("type %s struct {\n", os.Args[1])

    for _, tm := range tms {
        fmt.Printf("    %-30s%-20s%-20s\n", formatField(tm.Field), formatType(tm.Type), "`json:\"" + tm.Field + "\"`" );
    }

    fmt.Println("}")
}

func formatField(s string) string {
    bs := []byte(s)
    bs[0] = []byte(strings.ToUpper( string(bs[0]) ))[0]

    for idx := range bs {
        if bs[idx] == '_' {
            bs[idx + 1] = []byte(strings.ToUpper( string(bs[idx + 1]) ))[0]
        }
    }

    s = string(bs)
    s = strings.Replace(s, "_", "", -1)
    s = strings.Replace(s, "Id", "ID", -1)

    return s
}

func formatType(s string) string {
    if strings.Contains(s, "Array of ") {
        s = strings.Replace(s, "Array of ", "[]", -1)
    } else if strings.Contains(s, "number") {
        s = strings.Replace(s, " number", "", -1)
    }

    switch s {
    case "Integer":
        return "int64"
    case "Float":
        return "float64"
    case "String":
        return "string"
    case "True":
        return "bool"
    case "False":
        return "bool"
    case "Boolean":
        return "bool"
    default:
        return "*" + s
    }
}
