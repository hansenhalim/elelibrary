<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps<{ query?: string }>()
const router = useRouter()
const text = ref(props.query ?? '')

watch(
  () => props.query,
  (query) => {
    text.value = query ?? ''
  },
)

function submit() {
  const q = text.value.trim()
  if (q) {
    router.push({ name: 'search', query: { q } })
  }
}
</script>

<template>
  <form @submit.prevent="submit">
    <div class="field has-addons">
      <div class="control is-expanded">
        <input
          v-model="text"
          class="input"
          type="search"
          placeholder="Search by title, author, or keyword"
          aria-label="Search books"
        />
      </div>
      <div class="control">
        <button class="button is-primary" type="submit" :disabled="!text.trim()">Search</button>
      </div>
    </div>
  </form>
</template>
