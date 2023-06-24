from collections import defaultdict
from typing import List, Dict, Union

from numpy import ndarray
from scipy import spatial
from torch import Tensor

from schema import GroupMateScore, SimilarUserGroups, UserBackground
from sentence_transformers import SentenceTransformer


def get_cluster():
    """Get user information and create clusters.

    :return:
    """
    # Get all user with their background
    raise NotImplementedError


def create_background_mapping(users_background: List[UserBackground]) -> Dict[str, str]:
    """Get all individual with user background and create a universal mapping.

    :param users_background:
    :return:
    """
    user_background_mapping = {}
    for user_background in users_background:
        if user_background.background in user_background_mapping:
            continue
        user_background_mapping[user_background.background] = user_background.username

    return user_background_mapping


def generate_sentence_embeddings(
        user_background_mapping: Dict[str, str]
) -> Dict[str, Union[list[Tensor], ndarray, Tensor]]:
    """Generate a sentence embeddings.

    :param user_background_mapping: Dictionary with background as key and username as value
    :return:
    """
    model = SentenceTransformer('sentence-transformers/all-MiniLM-L6-v2')
    paragraphs = list(user_background_mapping.keys())
    embeddings = model.encode(paragraphs)

    user_embeddings = {
        user_background_mapping[
            paragraph]: embedding for paragraph, embedding in zip(paragraphs, embeddings)
    }

    return user_embeddings


def create_clusters(
        user_background_mapping: Dict[str, str],
        user_embeddings: Dict[str, Union[list[Tensor], ndarray, Tensor]]
) -> List[SimilarUserGroups]:
    """Create cluster and return dictionary with user and their top matches.

    :param user_background_mapping:
    :param user_embeddings:
    :return:
    """
    user_groups = []
    all_personal_pairs = defaultdict(list)
    for user in user_background_mapping.values():
        for user1 in user_background_mapping.values():
            all_personal_pairs[user].append(
                [spatial.distance.cosine(user_embeddings[user1], user_embeddings[user]), user1]
            )

    for user in user_background_mapping.values():
        groupmates = []
        for pair in sorted(all_personal_pairs[user], key=lambda x: x[1]):
            groupmates.append(GroupMateScore(
                username=pair[1],
                similarity_score=pair[0]
            ))
        SimilarUserGroups(
            username=user,
            groupmates=groupmates
        )
        user_groups.append(SimilarUserGroups)

    return user_groups
