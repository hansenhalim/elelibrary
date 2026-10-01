<script setup lang="ts">
import raterJs, { type Rater } from 'rater-js'
import { onBeforeUnmount, onMounted, useTemplateRef, watch } from 'vue'

const props = defineProps<{ rating: number }>()
const el = useTemplateRef<HTMLElement>('el')
let rater: Rater | undefined

onMounted(() => {
  if (el.value) {
    rater = raterJs({ element: el.value, rating: props.rating, readOnly: true, starSize: 20 })
  }
})

watch(
  () => props.rating,
  (rating) => rater?.setRating(rating),
)

onBeforeUnmount(() => rater?.dispose())
</script>

<template>
  <div ref="el"></div>
</template>
