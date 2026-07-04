"""
Vector store for DormGo agent: BGE embedding + FAISS ANN search.
Replaces SQL LIKE keyword search with semantic vector search.

Builds index from Go backend's /internal/agent/tools/all_posts endpoint,
refreshes periodically in background.
"""

import logging
import threading
import time
from typing import Optional

import faiss
import numpy as np

logger = logging.getLogger("dormgo.vector_store")

# Lazy imports — model loads on first use
_embedding_model: Optional[object] = None
_index: Optional[faiss.Index] = None
_post_ids: list[int] = []
_lock = threading.Lock()

# Model settings
MODEL_NAME = "BAAI/bge-small-zh-v1.5"
VECTOR_DIM = 512


def _get_model():
    """Lazy-load the embedding model (downloads ~100MB on first run).
    Uses local_files_only to avoid HuggingFace network timeouts in China."""
    global _embedding_model
    if _embedding_model is None:
        from sentence_transformers import SentenceTransformer
        logger.info(f"Loading embedding model: {MODEL_NAME} (offline mode) ...")
        _embedding_model = SentenceTransformer(MODEL_NAME, local_files_only=True)
        logger.info(f"Model loaded, dim={_embedding_model.get_sentence_embedding_dimension()}")
    return _embedding_model


def _embed_texts(texts: list[str]) -> np.ndarray:
    """Embed a list of texts into 512-dim vectors, normalized for cosine similarity."""
    model = _get_model()
    embeddings = model.encode(
        texts,
        normalize_embeddings=True,  # enables inner-product = cosine similarity
        show_progress_bar=False,
        batch_size=32,
    )
    return np.array(embeddings, dtype=np.float32)


def build_index(posts: list[dict], limit: int = None) -> int:
    """
    Build FAISS index from a list of post dicts.
    Each post must have 'id', 'title', 'content'.
    Returns the number of posts indexed.
    """
    global _index, _post_ids

    if limit:
        posts = posts[:limit]

    texts = [(p.get("title") or "") + " " + (p.get("content") or "") for p in posts]
    ids = []
    for p in posts:
        pid = p.get("id") or p.get("ID")
        try:
            ids.append(int(pid))
        except (TypeError, ValueError):
            ids.append(0)

    if not texts:
        logger.warning("No posts to index")
        with _lock:
            _index = faiss.IndexFlatIP(VECTOR_DIM)
            _post_ids = []
        return 0

    vectors = _embed_texts(texts)

    idx = faiss.IndexFlatIP(VECTOR_DIM)
    idx.add(vectors)

    with _lock:
        _index = idx
        _post_ids = ids

    logger.info(f"FAISS index built: {idx.ntotal} posts, dim={VECTOR_DIM}")
    return idx.ntotal


def search(query: str, k: int = 20) -> list[dict]:
    """
    Semantic search: embed query → FAISS ANN → return top-k matches.
    Returns list of {"id": post_id, "score": cosine_similarity}.
    """
    with _lock:
        if _index is None or _index.ntotal == 0:
            return []
        idx = _index
        ids = list(_post_ids)

    query_vec = _embed_texts([query])
    scores, indices = idx.search(query_vec, min(k, idx.ntotal))

    results = []
    for score, pos in zip(scores[0], indices[0]):
        if pos < 0 or pos >= len(ids):
            continue
        results.append({"id": ids[pos], "score": float(score)})
    return results


def is_ready() -> bool:
    """Check if the vector index has been built and is searchable."""
    with _lock:
        return _index is not None and _index.ntotal > 0


def index_size() -> int:
    with _lock:
        return _index.ntotal if _index else 0
