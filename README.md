# tgdocs2go
Simple tool to parse telegram documentation and generate go code.

For example, in order to generate a struct for the [User][1] type, one needs to run 

```
tgdocs2go User
```

and the following will be printed on stdout

```
// https://core.telegram.org/bots/api#User
type User struct {
    ID                            int                 `json:"id"`
    FirstName                     string              `json:"first_name"`
    LastName                      string              `json:"last_name"`
    Username                      string              `json:"username"`
    LanguageCode                  string              `json:"language_code"`
}
```

[1]: https://core.telegram.org/bots/api#user
