# [VMail](https://vmail.leopard.in.ua/) - check HTML & CSS compatibility with email clients [![Build and Deploy](https://github.com/le0pard/vmail/actions/workflows/deploy.yml/badge.svg?branch=main)](https://github.com/le0pard/vmail/actions/workflows/deploy.yml)

Email clients use different rendering standards. This is why your email can be displayed not as you designed it. You need to check that your message code won't cause rendering issues.

Vmail (**V**alidate E**mail**) check HTML and CSS email content in search of problematic elements. For each it finds, it displays the list of email clients that lack support for it or support it only partially.

## Development

Web app build on top of [middleman](http://middlemanapp.com/). To start it in development mode, you need install ruby, node.js, golang and run in terminal:

```bash
$ bundle # get all ruby deps
$ yarn # get all node.js deps
$ bundle exec rake wasm:compile # build wasm module
$ bundle exec middleman server # start server on 4567 port
```

## Build wasm file from Go files

```bash
bundle exec rake wasm:compile
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
