import os
from python_service.summarizers.summarizer_base import SummarizerBase
from transformers import TFAutoModelForSeq2SeqLM, AutoTokenizer, pipeline

MODEL_CACHE_PATH = "../../model_cache"

class HuggingFaceSummarizer(SummarizerBase):
    def __init__(self, model_name):
        super().__init__(model_name)
        model_name = self.model_configs["model_name"]
        tokenizer_name = self.model_configs["tokenizer_name"]

        if not os.path.exists(f"{MODEL_CACHE_PATH}/{model_name}/{model_name}.h5"):
            self.model = TFAutoModelForSeq2SeqLM.from_pretrained(model_name)
            self.model.save_pretrained(f"{MODEL_CACHE_PATH}/{model_name}")
            self.tokenizer = AutoTokenizer.from_pretrained(tokenizer_name)
            self.tokenizer.save_pretrained(f"{MODEL_CACHE_PATH}/{model_name}")
        else:
            self.model = TFAutoModelForSeq2SeqLM.from_pretrained(f"{MODEL_CACHE_PATH}/{model_name}")
            self.tokenizer = AutoTokenizer.from_pretrained(f"{MODEL_CACHE_PATH}/{model_name}")

    def summarize(self, text):
        input_ids = self.tokenizer.encode(text, return_tensors='tf')
        ids = self.model.generate(input_ids, max_length=self.model_configs["max_length"],
                            min_length=self.model_configs["min_length"])
        return self.tokenizer.decode(ids[0], skip_special_tokens=True)

if __name__ == "__main__":
    with open(file="/home/shalin/go_projects/reddit_posts_summarizer/reddit_post_summarizer/comments_344.txt") as f:
        comments = f.read()

    print(HuggingFaceSummarizer("t5-small").summarize(comments))