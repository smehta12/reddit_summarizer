# Reddit Thread Summarizer

This project summarizes reddit thread. It provides the full solution for the summarization that how it can be handeled for the fast responses. It uses OpenAPI Text Da-Vinci model and 
Huggingface T5 models to summarize. In the end it chooses best one using unsupervised 
evaluation methods among multiple models and returns the text to user.

This project has 2 parts. 
### Go Based
It provides REST API to initiate the request where user provides 
Reddit credentials and link to the post the user wants to summarize. It also embeds
interaction with the OpenAPI model to get the text summary. It uses following techniques
to enable summarizer.
   1) **Gin**: For web interface.
   2) **Reddit REST API**: For getting Bearer user token from Reddit and comments data from the thread.
   3) **OpenAPI REST API**: For getting summary from the text-davinci-003 model and 
   Cleanup of summary using  text-davinci-edit model.
   4) **Custom Python API**: Used for getting summary from the huggingface T5 transformer model.

   It creates recursive calls for summarization if the text is bigger than model's max
   size inputs. 
   It uses **_Go Interfaces_** to prepare recursive summarization of all types of 
   model. It also uses **_channels_, _go routines_** for parallel summarization from the different models.

### Python Based
It provides REST API for getting summarization from the different 
flavors of the huggingface T5 models. Uses following techniques to enable this.
   1) **FastAPI**: For providing REST API interface
   2) **Huggingface API**: To downloads the models and generate summary

### Unsupervised Summary Ranking
It is a microservice created using [SUPERT](https://arxiv.org/abs/2005.03724). The FastAPI is used for the REST API interface. It takes Reddit comments as a doc and generated summary as an inputs and ranks the summaries to find the best. SUPERT rates the quality of a summary by measuring its semantic similarity with a pseudo reference summary, i.e. selected salient sentences from the source documents, using contextualized embeddings and soft token alignment techniques.

### TODO:
Apply celery+RabbitMQ type solution for improved task synchronization. Currently it the summarization chuncks hangs the python service.
