<script setup lang="ts">
import { ref } from 'vue'
import { setFavorite } from '@/api/books'

const props = defineProps<{ bookId: string }>()
const favorite = defineModel<boolean>({ required: true })
const saving = ref(false)
const failed = ref(false)

async function toggle() {
  const next = !favorite.value
  saving.value = true
  failed.value = false
  try {
    await setFavorite(props.bookId, next)
    favorite.value = next
  } catch (err) {
    console.error(err)
    failed.value = true
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div>
    <button
      class="button is-small"
      :class="{ 'is-loading': saving }"
      type="button"
      :disabled="saving"
      @click="toggle"
    >
      {{ favorite ? '★ Favorited' : '☆ Favorite' }}
    </button>
    <p v-if="failed" class="help is-danger">Couldn't update favorite.</p>
  </div>
</template>
