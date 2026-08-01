import { afterEach, vi } from 'vitest'
import { config } from '@vue/test-utils'

afterEach(() => {
  document.body.innerHTML = ''
})

config.global.stubs = {
  RouterLink: {
    props: ['to'],
    template: '<a :href="typeof to === \'string\' ? to : \'#\'"><slot /></a>',
  },
}

vi.stubGlobal('ResizeObserver', class {
  observe() {}
  unobserve() {}
  disconnect() {}
})

window.scrollTo = vi.fn()
