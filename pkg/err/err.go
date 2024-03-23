package err

func CheckErr(err error, errFun func(error)) {
	if err != nil {
		errFun(err)
	}
}
