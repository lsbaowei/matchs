# matchs use

var dict = []string{"xxx", "x1"}  

//match keywords  

mc := matchs.NewMatchService()  

mc.Build(dict)  

words, _ := mc.Match(contents, false, '*')  



