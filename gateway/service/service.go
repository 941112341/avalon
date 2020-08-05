package service

type Gateway interface {
	Upload()   // idl 文件上传
	Transfer() // 转发
}

type UploadVoVo struct {
	PSM      string
	Filename string
	Body     string
}

type UploadService interface {
	Upload(request *UploadVoVo) error
	Get(request *UploadVoVo) (*UploadVoVo, error)
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
