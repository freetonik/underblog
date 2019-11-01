[![Build Status](https://travis-ci.org/freetonik/underblog.svg?branch=master)](https://travis-ci.org/freetonik/underblog)

# Underblog

An extremely simple, fast static blog generator.

## Install

On MacOS:

```
brew install freetonik/tap/underblog
```

Docker:

```
docker run --rm -it -v /path/to/your/blog:/blog freetonik/underblog
```

Other platforms: coming soon. Or you can build yourself by cloning the repo and running `make build`.

## How it works

You only need 4 things:

1. `index.html` template for blog's index page.
2. `post.html` template for single post.
3. `css/styles.css` for CSS styles.
3. `markdown` folder.

There is no front-matter. **Date** and **slug** are derived from the filename. **Title** is derived from the first line of markdown file. Make sure the first line starts with `#`.

**Step 1:** create the following folder structure:

```
.
├── css
│   └── styles.css
├── markdown
│   └── YYYY-MM-DD-Slug_1.md
│   └── YYYY-MM-DD-Slug_2.md
│   └── YYYY-MM-DD-Slug_3.md
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
- NO plugins
- NO dependencies

## Roadmap

- [x] derive dates from filenames
- [ ] RSS generation
- [ ] Syntax highlighting for code with Chroma
- [ ] live preview server (?)
