from abc import abstractmethod
import yaml
from yaml import Loader

class SummarizerBase(object):
    def __init__(self):
        # Get model configs
        with open("../../common/model_configs.yaml") as f:
            self._model_configs = yaml.load(f, Loader)

    @property
    def model_configs(self):
        return self._model_configs

    @abstractmethod
    def summarize(self, text):
        raise NotImplementedError()


# simple factory method to find model
def get_model_object(model_name):
    if model_name == 't5_small':
        from python_service.summarizers.huggingface import HuggingFaceSummarizer
        return HuggingFaceSummarizer(model_name)