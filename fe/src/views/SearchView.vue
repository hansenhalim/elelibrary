<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { searchBooks, type Book } from '@/api/books'
import SearchForm from '@/components/SearchForm.vue'
import StarRating from '@/components/StarRating.vue'

const PAGE_SIZE = 20

const route = useRoute()
const router = useRouter()
const q = computed(() => (typeof route.query.q === 'string' ? route.query.q.trim() : ''))

const books = ref<Book[]>([])
const loading = ref(false)
const failed = ref(false)
const done = ref(false)
let offset = 0

let requestId = 0

async function loadMore() {
  const id = ++requestId
  loading.value = true
  failed.value = false
  try {
    const page = await searchBooks(q.value, offset, PAGE_SIZE)
    if (id !== requestId) return

    const seen = new Set(books.value.map((book) => book.id))
    for (const book of page) {
      if (!seen.has(book.id)) {
        seen.add(book.id)
        books.value.push(book)
      }
    }
    offset += PAGE_SIZE
    done.value = page.length === 0
  } catch (err) {
    if (id !== requestId) return
    console.error(err)
    failed.value = true
  } finally {
    if (id === requestId) loading.value = false
  }
}

watch(
  q,
  (query) => {
    if (!query) {
      router.replace({ name: 'home' })
      return
    }
    books.value = []
    offset = 0
    done.value = false
    loadMore()
  },
  { immediate: true },
)
</script>

<template>
  <header class="section py-4 has-background-light">
    <div class="container">
      <div class="columns is-vcentered">
        <div class="column is-narrow">
          <RouterLink :to="{ name: 'home' }" class="title is-4">Elelibrary</RouterLink>
        </div>
        <div class="column">
          <SearchForm :query="q" />
        </div>
      </div>
    </div>
  </header>

  <main class="section">
    <div class="container">
      <article v-for="book in books" :key="book.id" class="media">
        <figure class="media-left cover">
          <img v-if="book.thumbnailUrl" :src="book.thumbnailUrl" :alt="book.title" loading="lazy" />
          <span v-else class="is-size-7 has-text-grey">No cover</span>
        </figure>
        <div class="media-content">
          <p class="has-text-weight-semibold">{{ book.title }}</p>
          <p class="has-text-grey mb-1">{{ book.authors.join(', ') || 'Unknown author' }}</p>
          <StarRating :rating="book.rating" />
        </div>
      </article>

      <p v-if="loading" class="mt-5">Loading…</p>
      <template v-else>
        <p v-if="failed" class="mt-5 has-text-danger">Couldn't load books. Please try again.</p>
        <p v-if="done" class="mt-5">
          {{ books.length ? 'No more results.' : `No books found for "${q}".` }}
        </p>
        <button v-else-if="books.length" class="button mt-5" type="button" @click="loadMore">
          Load more
        </button>
      </template>
    </div>
  </main>
</template>

<style scoped>
.cover {
  width: 64px;
  height: 96px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
}

.cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
