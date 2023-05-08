from python_service.summarizers.summarizer_base import SummarizerBase
from transformers import pipeline


class HuggingFaceSummarizer(SummarizerBase):
    def __init__(self, model_name):
        super().__init__()
        self.model_name = model_name

    def summarize(self, text):
        model_configs = self.model_configs[self.model_name]
        summarizer = pipeline("summarization", model=model_configs["model_name"],
                              tokenizer=model_configs["tokenizer_name"])
        return summarizer(text, max_length=model_configs["max_length"], min_length=model_configs["min_length"])


if __name__ == "__main__":
    with open(file="/home/shalin/go_projects/reddit_posts_summarizer/reddit_post_summarizer/comments_344.txt") as f:
        comments = f.read()

    print(HuggingFaceSummarizer("t5-small").summarize(comments))