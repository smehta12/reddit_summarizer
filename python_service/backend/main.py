import uvicorn

from fastapi import FastAPI
from pydantic import validator, BaseModel
from python_service.summarizers.summarizer_base import get_model_object, SummarizerBase

class SummarizerInput(BaseModel):
    model_name: str
    prompt: str

    @validator('model_name', allow_reuse=True)
    def validator_model_name(cls, model_name):
        model_names = SummarizerBase.get_model_configs().keys()
        if model_name not in model_names:
            raise ValueError(f"model_name must be one of {model_names}")

        return model_name


app = FastAPI()

@app.get("/", tags=["root"])
def home():
    return "Welcome to Summarization Service!"

@app.post("/summarize", tags=["summarization"])
async def summarize(summarizer_input: SummarizerInput):
    model = get_model_object(summarizer_input.model_name)
    return model.summarize(summarizer_input.prompt)


@app.get("/model_names", tags=["summarization"])
def get_model_names():
    return list(SummarizerBase.get_model_configs().keys())


if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True, log_level="info")