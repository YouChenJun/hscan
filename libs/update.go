package libs

type Update struct {
	VulnUpdate    bool
	LocalMetaData string //本地元数据地址
	GenerateMeta  string //生成元数据的地址
	MetadataUlr   string //元数据地址
	GitRepoUrl    string //git仓库地址
}
