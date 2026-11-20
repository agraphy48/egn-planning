import type { Directive, DirectiveBinding } from 'vue'
import { useAuthStore } from '../stores/authStore'

function checkPermission(el: HTMLElement, binding: DirectiveBinding<string | string[]>) {
  const authStore = useAuthStore()
  const { value } = binding

  if (!value) return

  let hasPermission = false

  if (Array.isArray(value)) {
    hasPermission = value.some((perm) => authStore.can(perm))
  } else if (typeof value === 'string') {
    hasPermission = authStore.can(value)
  }

  if (!hasPermission) {
    // Hide element cleanly
    el.style.display = 'none'
  } else {
    // Restore display if previously hidden
    if (el.style.display === 'none') {
      el.style.removeProperty('display')
    }
  }
}

export const canDirective: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string | string[]>) {
    checkPermission(el, binding)
  },
  updated(el: HTMLElement, binding: DirectiveBinding<string | string[]>) {
    checkPermission(el, binding)
  },
}
