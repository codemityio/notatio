export SHELL:=/bin/bash

default: test

noop:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/noop --output-format=svg

test:
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/plantuml/one/two --output-format=svg
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/plantuml/one/two --output-format=png --recursive
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest plantuml --input-path=testdata/plantuml/text.puml --output-format=svg
	docker run --rm \
		-e INPUT_PATH=testdata/graphviz/input.dot \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest graphviz --output-format=svg
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest graphviz --input-path=. --output-format=svg --recursive
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest mermaid --input-path=testdata/mermaid/input.mmd
	docker run --rm -it \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest mermaid --input-path=testdata/mermaid/one --recursive
	docker run --rm -it \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest mermaid --input-path=testdata/mermaid/one --output-format=png --recursive
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest mermaid --input-path=testdata/mermaid/sequence.mmd
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest mermaid --input-path=testdata/mermaid/class.mmd
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest mermaid --input-path=testdata/mermaid/state.mmd
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest mermaid --input-path=testdata/mermaid/entity-relationship.mmd
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest mermaid --input-path=testdata/mermaid/pie.mmd
	docker run --rm \
		-w ${PWD} \
		-v ${PWD}:${PWD} \
		codemityio/notatio:latest mermaid --input-path=testdata/mermaid/flow-other.mmd
