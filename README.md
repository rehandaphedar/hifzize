# Introduction

A program to generate page wise Anki flashcards for the Qurʾān.

It is available on HuggingFace as [hifzize](https://huggingface.co/datasets/rehandaphedar/hifzize).

Due to deck size restrictions, it is not available on AnkiWeb yet.

# Installation

```sh
go install git.sr.ht/~rehandaphedar/hifzize@latest
```

# Usage

The documentation for usage and flags can be accessed by running `hifzize -h`.

- The `-words` data can be obtained from QUL's [Ayah by ayah and word by text of Quran](https://qul.tarteel.ai/resources/quran-script)
- The `-layout` data can be obtained from QUL's [Mushaf Layout Resources](https://qul.tarteel.ai/resources/mushaf-layout)
- The `-phrases` data can be obtained from QUL's [Mutashabihat ul Quran - mutashabihat(Phrase)](https://qul.tarteel.ai/resources/mutashabihat/73)
- The `-metadata-*` can be obtained from QUL's [Quran data, surahs, ayahs, words, juz etc.](https://qul.tarteel.ai/resources/quran-metadata)
- The `-media-config` is a YAML file with a list of objects with the keys `src` and `as`. The filepaths are resolved relative to the config file.

The page images for the default deck are extracted using `pdfimages` from [the KFGQPC website](https://qurancomplex.gov.sa/), specifically, the 604 pages version of Al-Muṣḥaf Al-Wasaṭ of the Ḥafṣ Qirāʿah. The reasons for this choice are:

- Ḥafṣ is the most popular Qirāʿah.
- Normally, Al-Muṣḥaf Al-ʿĀdī would be used for the default deck, however, due to a peculiarity in its PDF, the width of the page image (extracted using `pdfimages` as mentioned before) changes if there are juz/ḥizb markers on the side. Thus, Al-Muṣḥaf Al-Wasaṭ has been used instead.

You are free to use another Muṣḥaf and/or Qirāʿah to generate your own deck.

# Note Type

The hifzize note type produces the Page Recall card type.

The front of the card shows few lines of the previous page.

The back of the card shows:
- the last few lines of the previous page
- the current page in full
- the first few lines of the next page
- the mutashābihāt for the current page
- any notes that the user has added in the Notes field

Clicking on any verse/page opens it in the Tarteel app.

![Front](https://git.sr.ht/~rehandaphedar/hifzize/blob/main/assets/front.png)
![Back 1](https://git.sr.ht/~rehandaphedar/hifzize/blob/main/assets/back-1.png)
![Back 2](https://git.sr.ht/~rehandaphedar/hifzize/blob/main/assets/back-2.png)
![Back 3](https://git.sr.ht/~rehandaphedar/hifzize/blob/main/assets/back-3.png)
![Back 4](https://git.sr.ht/~rehandaphedar/hifzize/blob/main/assets/back-4.png)

# Recommended Usage

It is recommended to suspend all notes at first and unsuspend by tag as you memorise.

It is also recommended to use the [mayyize](https://sr.ht/~rehandaphedar/mayyize) deck alongside this deck, unsuspending by page tag in the mayyize deck as you memorise. Note that the default tag format of both the decks is the same to facilitate this.
