Miniflux 2
==========

NOTE:
Gatra Bali project uses Miniflux 2 as the core part of the system, Miniflux aggregates news from several Balinese online media. I added Google Cloud Platform PubSub integration to publish a message to a Topic every time Category, Feed and Entry is created/updated/deleted. Every message sent to that topic will be received by the Cloud Functions (please check [gatrabali-functions](https://github.com/apps4bali/gatrabali-functions)) then based on the message the Cloud Function will call Miniflux REST API and process/store the response to Firestore.

Then the mobile app just need to use Firebase SDK to pull the news and all other features (read later, sharing, auth, comments, etc.) will be implemented on the Firebase side.


[![Build Status](https://travis-ci.org/miniflux/miniflux.svg?branch=master)](https://travis-ci.org/miniflux/miniflux)
[![GoDoc](https://godoc.org/miniflux.app?status.svg)](https://godoc.org/miniflux.app)
[![Documentation Status](https://readthedocs.org/projects/miniflux/badge/?version=latest)](https://docs.miniflux.app/)

Miniflux is a minimalist and opinionated feed reader:

- Written in Go (Golang)
- Works only with Postgresql
- Doesn't use any ORM
- Doesn't use any complicated framework
- Use only modern vanilla Javascript (ES6 and Fetch API)
- Single binary compiled statically without dependency
- The number of features is voluntarily limited

It's simple, fast, lightweight and super easy to install.

Official website: <https://miniflux.app>

Documentation
-------------

The Miniflux documentation is available here: <https://docs.miniflux.app/> ([Man page](https://miniflux.app/miniflux.1.html))

- [Opinionated?](https://docs.miniflux.app/en/latest/opinionated.html)
- [Features](https://docs.miniflux.app/en/latest/features.html)
- [Requirements](https://docs.miniflux.app/en/latest/requirements.html)
- [Installation Instructions](https://docs.miniflux.app/en/latest/installation.html)
- [Installation Tutorials](https://docs.miniflux.app/en/latest/tutorials.html)
- [Upgrading to a New Version](https://docs.miniflux.app/en/latest/upgrade.html)
- [Configuration](https://docs.miniflux.app/en/latest/configuration.html)
- [Command Line Usage](https://docs.miniflux.app/en/latest/cli.html)
- [User Interface Usage](https://docs.miniflux.app/en/latest/usage.html)
- [Keyboard Shortcuts](https://docs.miniflux.app/en/latest/keyboard_shortcuts.html)
- [Integration with External Services](https://docs.miniflux.app/en/latest/integration.html)
- [Scraper Rules](https://docs.miniflux.app/en/latest/scraper_rules.html)
- [Rewrite Rules](https://docs.miniflux.app/en/latest/rewrite_rules.html)
- [REST API](https://docs.miniflux.app/en/latest/api.html)
- [Development](https://docs.miniflux.app/en/latest/development.html)
- [Internationalization](https://docs.miniflux.app/en/latest/i18n.html)
- [Frequently Asked Questions](https://docs.miniflux.app/en/latest/faq.html)

Screenshots
-----------

Default theme:

![Default theme](https://miniflux.app/image/overview.png)

Dark theme when using keyboard navigation:

![Dark theme](https://miniflux.app/image/item-selection-black-theme.png)

Credits
-------

- Authors: Frédéric Guillot - [List of contributors](https://github.com/miniflux/miniflux/graphs/contributors)
- Distributed under Apache 2.0 License
