// func TestMain_bind_address_exist(t *testing.T) {
// 	assert := assert.New(t)

// 	_, err := net.Listen("tcp", "0.0.0.0:8080")
// 	if err != nil {
// 		log.Error().Err(err).Msg("Error occured while creating server on port 8080")
// 	} else {
// 		log.Info().Msg("Listening on port 8080")
// 	}
// 	f1 := func() {
// 		main()
// 	}
// 	output := tools.CaptureOutput(f1)
// 	log.Info().Msgf("output %s", output)
// 	assert.Contains(output, "Startup failed")
// }

// func TestMain_bind_address_exist(t *testing.T) {
// 	_, err := net.Listen("tcp", "0.0.0.0:8080")
// 	if err != nil {
// 		log.Error().Err(err).Msg("Error occured while creating server on port 8080")
// 	} else {
// 		log.Info().Msg("Listening on port 8080")
// 	}
// 	proc, err := os.FindProcess(os.Getpid())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	sigc := make(chan os.Signal, 2)
// 	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
// 	go func() {
// 		main()
// 		sig := <-sigc
// 		switch sig {
// 		case os.Interrupt:
// 			log.Info().Msg("plop")
// 		case syscall.SIGTERM:
// 			log.Info().Msg("bbbbzzzz")
// 		}
// 		signal.Stop(sigc)
// 	}()

// 	proc.Signal(os.Interrupt)
// 	time.Sleep(2 * time.Second)
// }

// func TestMain_bind_address_exist(t *testing.T) {
// 	_, err := net.Listen("tcp", "0.0.0.0:8080")
// 	if err != nil {
// 		log.Error().Err(err).Msg("Error occured while creating server on port 8080")
// 	} else {
// 		log.Info().Msg("Listening on port 8080")
// 	}
// 	proc, err := os.FindProcess(os.Getpid())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	sigc := make(chan os.Signal, 1)
// 	signal.Notify(sigc, os.Interrupt)
// 	var output string

// 	go func(output string) {
// 		f1 := func() {
// 			main()
// 		}
// 		output = tools.CaptureOutput(f1)
// 		<-sigc
// 		signal.Stop(sigc)
// 	}(output)

// 	log.Info().Msgf("output %v", output)
// 	time.Sleep(1 * time.Second)
// 	proc.Signal(os.Interrupt)
// }

func TestMain_bind_address_exist(t *testing.T) {
	assert := assert.New(t)
	_, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Error().Err(err).Msg("Error occured while creating server on port 8080")
	} else {
		log.Info().Msg("Listening on port 8080")
	}
	_, err = os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(
		sigc,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	var output string

	go func(output string) {
		f1 := func() {
			main()
		}
		output = tools.CaptureOutput(f1)
		// main()
		<-sigc
		// log.Info().Msgf("sig %v", sig)
		signal.Stop(sigc)
	}(output)

	// log.Info().Msgf("output %v", output)
	// time.Sleep(1 * time.Second)
	// proc.Signal(os.Interrupt)
	assert.Contains(output, "listen tcp :8080: bind: address already in use")
	assert.Erro
}