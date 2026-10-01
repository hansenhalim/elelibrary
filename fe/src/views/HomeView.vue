<script setup lang="ts">
import SearchForm from '@/components/SearchForm.vue'
import logo from '@/assets/logo.svg'
</script>

<template>
  <section class="hero is-fullheight">
    <div class="hero-body">
      <div class="container">
        <div class="columns is-centered">
          <div class="column is-half">
            <h1 class="has-text-centered mb-5">
              <img :src="logo" alt="Elelibrary" style="max-width: 300px" />
            </h1>
            <SearchForm />
            <p class="help has-text-centered mt-3">
              Can't find good books, why don't give
              <RouterLink :to="{ name: 'search', query: { q: 'PHP books' } }">
                PHP books
              </RouterLink>
              a try :)
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
