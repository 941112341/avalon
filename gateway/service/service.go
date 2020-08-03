package service

type Gateway interface {
	Upload()   // idl 文件上传
	Transfer() // 转发
}

type IDLFileVo struct {
	PSM      string
	Filename string
	Body     string
}

type UploadService interface {
	Upload(request *IDLFileVo) error
	Get(request *IDLFileVo) (*IDLFileVo, error)
}

type MapperID struct {
	URL    string
	Domain string
}

type MappingsService interface {
	GetMappers()
	AddMappers()
	DelMappers()
}
