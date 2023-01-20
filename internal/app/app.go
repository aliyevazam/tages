package app


func Run(cfg config.Config) {
	l := logger.New(cfg.LogLevel)

	fileStore := filestore.NewDiskFileStore("files")

	TagesService := service.NewTagesService(l, fileStore)

	lis, err := net.Listen("tcp", ":"+cfg.ServicePort)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - grpcclient.New: %w", err))
	}

}