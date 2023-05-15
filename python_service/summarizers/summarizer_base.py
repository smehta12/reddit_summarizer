from abc import abstractmethod
import yaml
from yaml import Loader

class SummarizerBase(object):

    @staticmethod
    def get_model_configs():
        with open("../../common/model_configs.yaml") as f:
            return yaml.load(f, Loader)

    def __init__(self, model_name):
        self._model_configs = self.get_model_configs()[model_name]

    @property
    def model_configs(self):
        return self._model_configs

    @abstractmethod
    def summarize(self, text):
        raise NotImplementedError()

# simple factory method to find model
def get_model_object(model_name):
    if model_name == 't5-small':
        from python_service.summarizers.huggingface import HuggingFaceSummarizer
        return HuggingFaceSummarizer(model_name)