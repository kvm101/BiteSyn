from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from deep_translator import GoogleTranslator
from textblob import TextBlob
from langdetect import detect
from better_profanity import profanity


profanity.load_censor_words()

app = FastAPI(title="NLP Review Rating API")

class ReviewRequest(BaseModel):
    text: str


def is_profane(text: str) -> bool:
    try:
        return profanity.contains_profanity(text)
    except Exception:
        return False


def get_star_rating_auto(text: str) -> int:
    try:
        lang = detect(text)
        translated = GoogleTranslator(source=lang, target='en').translate(text) if lang != 'en' else text
        polarity = TextBlob(translated).sentiment.polarity
        rating = (polarity + 1) * 2
        return max(1, min(5, round(rating)+1))
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"NLP помилка: {str(e)}")


@app.post("/rate")
async def rate_review(request: ReviewRequest):
    original_text = request.text
    try:
        lang = detect(original_text)
        translated = GoogleTranslator(source=lang, target='en').translate(original_text) if lang != 'en' else original_text
    except Exception:
        translated = original_text


    if is_profane(translated):
        return {
            "review": original_text,
            "status": False,
            "rating": 0
        }

    rating = get_star_rating_auto(original_text)
    return {
        "review": original_text,
        "status": True,
        "rating": float(rating)
    }
