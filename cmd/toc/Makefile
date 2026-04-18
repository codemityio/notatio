SHELL := /bin/bash

%: # Silence errors about non existing targets
	@true

docs:
	@notatio coi --command="${BASE_NAME} toc --document-path=README.md --header='Table of contents' int --help" --document-path=README.md --header='`int`' --limiter-left=### --limiter-right="### " --index=1
	@notatio coi --command="${BASE_NAME} toc --document-path=README.md --header='Table of contents' ext --help" --document-path=README.md --header='`ext`' --limiter-left=### --limiter-right="## " --index=1
