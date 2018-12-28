# Better search engine

All of the modern e-commerce search engines have major problems. Retailers use "sort by sales volume" to cover up deficiencies. All kind of snake oil search personalization tools promise wonderful results but being based on group behavior, they tend to perpetuate same products over and over -- you might as well reduce your assortment to a 100 items.

That said, I'd like to build a decent and free catalog backed search engine from the ground up. And I want to build it in Go. And I want to play with machine learning as I "Go". 

End product is supposed to have the following qualities:

* Facet search, including hierarchical facets
* Machine learning based data enrichment and discovery
* Easy consumption of existing catalog data
* Higher accuracy of search results even in relevancy search mode
* Easy to use data ingest and API
* Kubernetes native

For now, I'm using a data set of 6 major categories from amazon.com to have the phrase training done. This might take a while, but I do have a plenty of time.