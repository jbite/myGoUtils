//InsertSep insert sep to string every n but no end
//example: n := InsertSep("12345678",":",2)
//get 12:34:56:78

func InsertSep(s string, sep string, n int) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if i%n == (n-1) && i != len(s)-1 {
			b.WriteByte(s[i])
			b.Write([]byte(sep))
		} else {
			b.WriteByte(s[i])
		}

	}
	o := b.String()
	return o
}
