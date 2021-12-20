# [VMail](https://vmail.leopard.in.ua/) - check the markup (HTML, CSS) of HTML email template compatibility with email clients

[![Build and Deploy](https://github.com/le0pard/vmail/actions/workflows/deploy.yml/badge.svg?branch=main)](https://github.com/le0pard/vmail/actions/workflows/deploy.yml)

[![VMail](https://user-images.githubusercontent.com/98444/142698496-ee804d5e-1108-47a0-95ba-6eedd72e7144.png)](https://vmail.leopard.in.ua/)

Email clients use different rendering standards. This is why your email can be displayed not as you designed it. You need to check that your message code won't cause rendering issues.

Vmail (**V**alidate E**mail**) check the markup (HTML, CSS) of HTML email template content in search of problematic elements. For each it finds, it displays the list of email clients that lack support for it or support it only partially.

VMail collect the data on support for particular HTML & CSS rules from [Caniemail.com](https://www.caniemail.com/)

## Development

Web app build on top of [middleman](http://middlemanapp.com/). To start it in development mode, you need install ruby, node.js, golang and run in terminal:

```bash
$ bundle # get all ruby deps
$ yarn # get all node.js deps
$ bundle exec rake wasm:parser # build wasm parser module
$ bundle exec rake wasm:inliner # build wasm inliner module
$ bundle exec middleman server # start server on 4567 port
```

### Build wasm files from Go files

```bash
$ bundle exec rake wasm:parser # build wasm parser module
$ bundle exec rake wasm:inliner # build wasm inliner module
```

### Format svelte components

```bash
yarn prettier --write --plugin-search-dir=. ./webpack/components/*
```

### Benchmark parser

```bash
$ cd wasm_parser/parser
$ go test -bench=. -benchmem
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

