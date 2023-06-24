from typing import List

import pydantic

from pydantic import BaseModel, Field


class UserBackground(BaseModel):
    username: str = Field(..., description='User name.')
    background: str = Field(..., description='Background of the user.')


class GroupMateScore(BaseModel):
    username: str = Field(..., description='User name.')
    similarity_score: float = Field(..., description='Similarity Score')


class SimilarUserGroups(BaseModel):
    username: str = Field(..., description='User name.')
    groupmates: List[GroupMateScore] = Field(..., description='List of groupmates.')
