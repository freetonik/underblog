# Underblog

## Install

On MacOS:

```
brew install freetonik/tap/underblog
```

Other platforms: coming soon. Or you can build yourself by cloning the repo and running `make build`.

## How it works

An extremely simple, fast static blog generator.

**Step 1:** create the following folder structure:

```
├── css
│   └── styles.css
├── markdown
│   └── DD-MM-YYYY-Slug_1.md
│   └── DD-MM-YYYY-Slug_2.md
│   └── DD-MM-YYYY-Slug_3.md
├── index.html
├── post.html
```

(See [/example](example))

**Step 2:** run `underblog`.

**Step 3:** Your site is generated in `public`.

## Features

- NO front matter
- NO themes
- NO JavaScript
- NO tags, categories, taxonomy
- NO template lookup logic

## Roadmap

- [x] derive dates from filenames
- [ ] RSS generation
- [ ] Syntax highlighting for code with Chroma
- [ ] live preview server (?)
