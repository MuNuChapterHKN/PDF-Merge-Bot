package bot

import (
    "os"
    "log"
    "io"
    "io/ioutil"
    "net/http"
    "strconv"
    "github.com/AntonioLangiu/pdf_merge_bot/common"
    "gopkg.in/telegram-bot-api.v4"
    unicommon "github.com/unidoc/unidoc/common"
    unipdf "github.com/unidoc/unidoc/pdf"
)

func LoadBot(configuration *common.Configuration) {
	// Start the Bot
	bot, err := tgbotapi.NewBotAPI(configuration.TelegramAPI)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
    if err != nil {
        log.Panic(err)
    }

    // Create pdf folder and init unidoc library
    initFolders("./files")
    // this logger is the dummy logger that does not log nothing. We should decide what to do with the logs of the library
    unicommon.SetLogger(unicommon.DummyLogger{})

	var out string
	for update := range updates {
        var base_dir string = "./files/"+strconv.FormatInt(update.Message.Chat.ID, 10)+"/"
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			var command string = update.Message.Command()
			// General Commands
			// If the message /start is received create a folder
            // for that specific user or group.
			// If it's already present print the help message
			// if help is received print the help message
			log.Printf("Command is %s\n\n",command)
            log.Printf("the base dir is %s\n\n", base_dir)

			if command == "start" {
				out = "This bot helps you to merge different pdf"
			} else if command == "help" {
				out = "This bot helps you to merge different pdf"
            /***** user commands ***********/
			} else if command == "init" {
                out = "Ok, I setted up everything for you. You can start adding files"
                initFolders(base_dir)
			} else if command == "merge" {
                if _, err := os.Stat(base_dir); os.IsNotExist(err) {
                    out = "You must first start with init and then call merge";
                } else {
                    files, err := ioutil.ReadDir(base_dir)
                    if err != nil {
                        // Log the error messsage and abort the operation.
                        log.Print(err)
                    }
                    pdfWriter := unipdf.NewPdfWriter()
                    for _, f := range files {
                        log.Println("adding "+base_dir+"/"+f.Name())
                        f, err := os.Open(base_dir+"/"+f.Name())
                        if err != nil {
                            log.Fatal(err)
                        }
                        defer f.Close()

                        pdfReader, err := unipdf.NewPdfReader(f)
                        if err != nil {
                            log.Fatal(err)
                        }
                        numPages, err := pdfReader.GetNumPages()
                        if err != nil {
                            log.Fatal(err)
                        }

                        for i := 0; i < numPages; i++ {
                            pageNum := i + 1
                            page, err := pdfReader.GetPage(pageNum)
                            if err != nil {
                                log.Fatal(err)
                            }
                            err = pdfWriter.AddPage(page)
                            if err != nil {
                                log.Fatal(err)
                            }
                        }
                    }
                    fWrite, err := os.Create(base_dir+"/merged.pdf")
                    if err != nil {
                        log.Fatal(err)
                    }
                    defer fWrite.Close()
                    err = pdfWriter.Write(fWrite)
                    if err != nil {
                        log.Fatal(err)
                    }
                    docConfig := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, base_dir+"/merged.pdf")
                    //var document tgbotapi.Document
                    //document.FileID = docConfig.BaseFile.FileID
                    bot.Send(docConfig)
                }
			}
		} else {
            if (update.Message.Document == nil) {
                out = "To add a file you need o send a file"
            } else {
                var file_name string = update.Message.Document.FileID
                file, err := bot.GetFileDirectURL(file_name);
                if err != nil {
                    log.Fatal(err)
                }
                pdf, err := os.Create(base_dir+file_name)
                if err != nil {
                    log.Fatal(err)
                }
                defer pdf.Close()
                response, err := http.Get(file)
                defer response.Body.Close()

                numBytesWritten, err := io.Copy(pdf, response.Body)
                if err != nil {
                    log.Fatal(err)
                }
                if numBytesWritten < 0 {
                    log.Fatal("numBytesWritten < 0")
                }
                out = "Your file has been correcly added"
            }
		}

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, out)
		bot.Send(msg)
	}
}

func initFolders(name string) {
    if _, err := os.Stat(name); ! os.IsNotExist(err) {
        os.RemoveAll(name)
    }
    err := os.Mkdir(name, 0777)
    if err != nil {
        log.Fatal(err)
    }
}
