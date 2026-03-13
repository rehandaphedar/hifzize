package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"git.sr.ht/~rehandaphedar/genanki-go-utils/v2/pkg/qul"
	"github.com/npcnixel/genanki-go"
)

func main() {
	mushafId := flag.String("mushaf-id", "hafs_wasat_604", "A string uniquely identifying the mushaf. Used for generating note ids.")
	images := flag.String("images", "images/hafs_wasat_604/%03d.png", "Location of image files. %d is replaced with the page number.")
	imagesInCollection := flag.String("images-in-collection", "hifzize-hafs_wasat_604-%03d.png", "Location of image files in the collection. %d is replaced with the page number.")
	specialPagesStr := flag.String("special-pages", "1,2", "Special pages (Will be rendered without mask even when not the current page. This behavior can be customised by editing the template.)")

	modelId := flag.Int64("model-id", int64(1109290091), "ID of the model")
	modelName := flag.String("model-name", "hifzize", "Name of the model")

	deckId := flag.Int64("deck-id", int64(1816673620), "ID of the deck")
	deckName := flag.String("deck-name", "hifzize", "Name of the peck")
	deckDescription := flag.String("deck-description", "Recall the current page and the beginning of the next page. Based on the 604 pages Ḥafṣ Wasaṭ Muṣḥaf.", "Description of the deck")

	outputPath := flag.String("output", "out/hifzize-hafs_wasat_604.apkg", "Output filepath")

	templateHtmlPath := flag.String("template-html", "templates/index.gohtml", "Path to template HTML file")
	templateCssPath := flag.String("template-css", "templates/style.css", "Path to template CSS file")
	templateQfmtPath := flag.String("template-qfmt", "templates/qfmt.html", "Path to template Qfmt file")
	templateAfmtPath := flag.String("template-afmt", "templates/afmt.html", "Path to template Afmt file")

	templateFrontName := flag.String("template-front", "front", "Name of the front template")
	templateBackName := flag.String("template-back", "back", "Name of the back template")

	wordsPath := flag.String("words", "data/qpc-hafs-word-by-word.json", "Path to words data")
	layoutPath := flag.String("layout", "data/qpc-v4-tajweed-15-lines.db", "Path to layout data")
	metadataAyahPath := flag.String("metadata-ayah", "data/quran-metadata-ayah.json", "Path to ayah metadata")
	metadataJuzPath := flag.String("metadata-juz", "data/quran-metadata-juz.json", "Path to juz metadata")
	metadataHizbPath := flag.String("metadata-hizb", "data/quran-metadata-hizb.json", "Path to hizb metadata")
	metadataRubPath := flag.String("metadata-rub", "data/quran-metadata-rub.json", "Path to rub metadata")
	metadataManzilPath := flag.String("metadata-manzil", "data/quran-metadata-manzil.json", "Path to manzil metadata")
	metadataRukuPath := flag.String("metadata-ruku", "data/quran-metadata-ruku.json", "Path to ruku metadata")

	var tagFormat qul.TagFormat

	tagFormat.Chapter = flag.String("tag-format-chapter", "quran::chapter::%03d", "Format of the chapter tag. %d is replaced with the chapter number.")
	tagFormat.Verse = flag.String("tag-format-verse", "quran::verse::%s", "Format of the verse tag. %s is replaced with the zero padded verse key (Example: 001:001).")
	tagFormat.Page = flag.String("tag-format-page", "quran::page::%03d", "Format of the page tag. %d is replaced with the page number.")
	tagFormat.Juz = flag.String("tag-format-juz", "quran::juz::%02d", "Format of the juz tag. %d is replaced with the juz number.")
	tagFormat.Hizb = flag.String("tag-format-hizb", "quran::hizb::%02d", "Format of the hizb tag. %d is replaced with the hizb number.")
	tagFormat.Rub = flag.String("tag-format-rub", "quran::rub::%03d", "Format of the rub tag. %d is replaced with the rub number.")
	tagFormat.Manzil = flag.String("tag-format-manzil", "quran::manzil::%d", "Format of the manzil tag. %d is replaced with the manzil number.")
	tagFormat.Ruku = flag.String("tag-format-ruku", "quran::ruku::%03d", "Format of the ruku tag. %d is replaced with the ruku number.")

	flag.Parse()

	var words map[string]qul.Word
	var metadataAyah map[string]qul.MetadataAyah

	var metadataDivision qul.MetadataDivision

	err := loadJSON(*wordsPath, &words)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataAyahPath, &metadataAyah)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataJuzPath, &metadataDivision.Juz)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataHizbPath, &metadataDivision.Hizb)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataRubPath, &metadataDivision.Rub)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataManzilPath, &metadataDivision.Manzil)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataRukuPath, &metadataDivision.Ruku)
	if err != nil {
		log.Fatal(err)
	}

	index, err := qul.BuildIndex(*layoutPath, words, metadataDivision, tagFormat)
	if err != nil {
		log.Fatalf("build index: %v", err)
	}

	qfmt, err := readFile(*templateQfmtPath)
	if err != nil {
		log.Fatal(err)
	}
	afmt, err := readFile(*templateAfmtPath)
	if err != nil {
		log.Fatal(err)
	}

	css, err := readFile(*templateCssPath)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles(*templateHtmlPath)
	if err != nil {
		log.Fatalf("parse template files: %v", err)
	}
	var buf bytes.Buffer

	model := genanki.NewModel(*modelId, *modelName).
		SetCSS(css).
		AddField(genanki.Field{Name: "Page"}).
		AddField(genanki.Field{Name: "Front"}).
		AddField(genanki.Field{Name: "Back"}).
		AddField(genanki.Field{Name: "Notes"}).
		AddTemplate(genanki.Template{
			Name: "Page Recall",
			Qfmt: qfmt,
			Afmt: afmt,
		})
	deck := genanki.NewDeck(*deckId, *deckName, *deckDescription)

	var specialPages []int
	for specialPageStr := range strings.SplitSeq(*specialPagesStr, ",") {
		specialPage, err := strconv.Atoi(strings.TrimSpace(specialPageStr))
		if err != nil {
			log.Fatalf("parse special page %s: %v", specialPageStr, err)
		}
		specialPages = append(specialPages, specialPage)
	}

	totalPages := 0
	for pageNumber := range index.Tag.Page {
		if pageNumber > totalPages {
			totalPages = pageNumber
		}
	}

	mediaFiles := map[string]string{}

	for pageNumber, pageTags := range index.Tag.Page {
		current := Page{
			Type:   PageTypeNormal,
			Number: pageNumber,
		}
		previous := Page{
			Type:   PageTypeNormal,
			Number: pageNumber - 1,
		}
		next := Page{
			Type:   PageTypeNormal,
			Number: pageNumber + 1,
		}

		if previous.Number == 0 {
			previous.Type = PageTypeOpening
		}
		if slices.Contains(specialPages, previous.Number) {
			previous.Type = PageTypeSpecial
		}
		if slices.Contains(specialPages, next.Number) {
			next.Type = PageTypeSpecial
		}
		if next.Number == totalPages+1 {
			next.Type = PageTypeConclusion
		}

		imageFilePathCurrent := fmt.Sprintf(*images, current.Number)
		imageFilePathInCollectionCurrent := fmt.Sprintf(*imagesInCollection, current.Number)
		current.Path = imageFilePathInCollectionCurrent

		mediaFiles[imageFilePathCurrent] = imageFilePathInCollectionCurrent

		if (previous.Type == PageTypeNormal) || (previous.Type == PageTypeSpecial) {
			imageFilePathInCollectionPrevious := fmt.Sprintf(*imagesInCollection, previous.Number)
			previous.Path = imageFilePathInCollectionPrevious
		}

		if (next.Type == PageTypeNormal) || (next.Type == PageTypeSpecial) {
			imageFilePathInCollectionNext := fmt.Sprintf(*imagesInCollection, next.Number)
			next.Path = imageFilePathInCollectionNext
		}

		templateData := TemplateData{
			Previous: previous,
			Current:  current,
			Next:     next,
		}

		templateErrorMessage := "execute template %s with data %+v: %v"

		err = tmpl.ExecuteTemplate(&buf, *templateFrontName, templateData)
		if err != nil {
			log.Printf(templateErrorMessage, *templateFrontName, templateData, err)
		}
		front := buf.String()
		buf.Reset()

		err = tmpl.ExecuteTemplate(&buf, *templateBackName, templateData)
		if err != nil {
			log.Printf(templateErrorMessage, *templateBackName, templateData, err)
		}
		back := buf.String()
		buf.Reset()

		note := genanki.NewNote(
			model.ID,
			[]string{
				strconv.Itoa(pageNumber),
				front,
				back,
				"",
			},
			pageTags,
		)

		noteIdBase := fmt.Sprintf("%d_%s_%d", model.ID, *mushafId, pageNumber)
		note.ID = qul.GenerateID(noteIdBase)
		deck.AddNote(note)
	}

	pkg := genanki.NewPackage([]*genanki.Deck{deck}).AddModel(model)

	for src, as := range mediaFiles {
		mediaFileData, err := os.ReadFile(src)
		if err != nil {
			log.Fatalf("read media file %s: %v", src, err)
		}
		pkg.AddMedia(as, mediaFileData)
	}

	if err := pkg.WriteToFile(*outputPath); err != nil {
		log.Fatalf("write package to %s: %v", *outputPath, err)
	}
}
