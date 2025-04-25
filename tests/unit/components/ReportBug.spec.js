import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import ReportBug from '@/pages/ReportBug.vue'
import { useBugStore } from '@/stores/bug'

describe('ReportBug.vue', () => {
  const createWrapper = () => {
    return mount(ReportBug, {
      global: {
        plugins: [createTestingPinia()]
      }
    })
  }

  test('validates required fields', async () => {
    const wrapper = createWrapper()
    await wrapper.find('form').trigger('submit')

    expect(wrapper.find('[data-test="title-error"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="description-error"]').exists()).toBe(true)
  })

  test('submits form with valid data', async () => {
    const wrapper = createWrapper()
    const store = useBugStore()

    await wrapper.find('[data-test="title-input"]').setValue('Test Bug')
    await wrapper.find('[data-test="description-input"]').setValue('Test Description')
    await wrapper.find('[data-test="priority-select"]').setValue('high')
    await wrapper.find('form').trigger('submit')

    expect(store.createBug).toHaveBeenCalledWith({
      title: 'Test Bug',
      description: 'Test Description',
      priority: 'high'
    })
  })

  test('shows loading state during submission', async () => {
    const wrapper = createWrapper()
    const store = useBugStore()
    store.createBug.mockImplementationOnce(() => new Promise(() => {}))

    await wrapper.find('[data-test="title-input"]').setValue('Test Bug')
    await wrapper.find('[data-test="description-input"]').setValue('Test Description')
    await wrapper.find('form').trigger('submit')

    expect(wrapper.find('[data-test="submit-button"]').attributes('disabled')).toBeTruthy()
    expect(wrapper.find('[data-test="submit-button"]').text()).toBe('Submitting...')
  })
}) 