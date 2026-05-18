<script setup langts>
import { ref } from 'vue'
import { useUiStore } from './stores/ui'

const uiStore = useUiStore()

// Dữ liệu giả lập mô phỏng chính xác theo hình ảnh thiết kế
const conversations = ref([
  { id: 1, title: 'This is the selected conversation', active: true },
  { id: 2, title: 'This is a conversation', active: false },
  { id: 3, title: 'This is a conversation', active: false },
  { id: 4, title: 'This is a conversation', active: false }
])

const messages = ref([
  { id: 1, role: 'user', text: ['This is the user text / prompt.', 'This is the user text / prompt.'] },
  { id: 2, role: 'assistant', text: ['This is the AI Response.', 'This is the AI Response.'] },
  { id: 3, role: 'user', text: ['This is the user text / prompt.', 'This is the user text / prompt.', 'This is the user text / prompt.'] },
  { id: 4, role: 'assistant', text: ['This is the AI Response.', 'This is the AI Response.', 'This is the AI Response.', 'This is the AI Response.'] },
  { id: 5, role: 'user', text: ['This is the user text / prompt.', 'This is the user text / prompt.'] },
  { id: 6, role: 'assistant', text: ['This is the AI Response.'] },
  { id: 7, role: 'user', text: ['This is the user text / prompt.'] },
  { id: 8, role: 'assistant', text: ['This is the AI Response.', 'This is the AI Response.'] }
])

const currentModel = ref('Gemma-4-26b-a4b-it')
const inputMessage = ref('')
</script>

<template>
  <div class="flex h-screen w-screen overflow-hidden bg-[#0d1113] text-white antialiased">
    
    <aside 
      :class="[
        uiStore.isSidebarCollapsed ? 'w-[72px]' : 'w-72',
        'h-full bg-[#161a1d] border-r border-[#232d30] flex flex-col justify-between transition-all duration-300 ease-in-out select-none'
      ]"
    >
      <div class="flex flex-col items-center w-full">
        <div class="h-20 w-full flex items-center px-5">
          <button 
            @click="uiStore.toggleSidebar()"
            class="w-10 h-10 flex items-center justify-center rounded-xl border border-[#2d3a3e] bg-[#1a2326] text-[#4ea399] hover:bg-[#233034] transition-colors cursor-pointer"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
        </div>

        <div class="w-full px-4 mb-6">
          <button 
            v-if="!uiStore.isSidebarCollapsed"
            class="w-full h-11 bg-[#37736b] hover:bg-[#2d5f58] text-white font-medium rounded-xl flex items-center justify-center gap-2 transition-all cursor-pointer truncate"
          >
            <span>+ New Conversation</span>
          </button>
          <button 
            v-else
            class="w-10 h-10 bg-[#37736b] hover:bg-[#2d5f58] text-white rounded-xl flex items-center justify-center transition-all cursor-pointer mx-auto"
          >
            <span class="text-xl font-light">+</span>
          </button>
        </div>

        <div class="w-full flex flex-col gap-1 px-2 overflow-y-auto max-h-[calc(100vh-250px)]">
          <p v-if="!uiStore.isSidebarCollapsed" class="text-xs font-bold text-[#8fa3a6] px-3 py-2 uppercase tracking-wider">
            Conversation
          </p>
          
          <div 
            v-for="item in conversations" 
            :key="item.id"
            :class="[
              item.active 
                ? 'bg-[#37736b] text-white' 
                : 'text-[#a6b8ba] hover:bg-[#1f2629] hover:text-white',
              uiStore.isSidebarCollapsed ? 'justify-center py-3' : 'px-3 py-2.5',
              'rounded-lg text-sm font-medium transition-all cursor-pointer flex items-center'
            ]"
          >
            <span v-if="!uiStore.isSidebarCollapsed" class="truncate">{{ item.title }}</span>
            <span v-else :class="[item.active ? 'text-white' : 'text-[#4ea399]', 'font-bold text-base']">■</span>
          </div>
        </div>
      </div>

      <div class="p-4 border-t border-[#232d30] text-center truncate">
        <span v-if="!uiStore.isSidebarCollapsed" class="text-xs font-semibold text-[#62777a]">
          Made by <span class="text-[#8fa3a6]">OiHuecuata</span>
        </span>
        <span v-else class="text-xs font-bold text-[#4ea399]">OH</span>
      </div>
    </aside>

    <main class="flex-1 h-full flex flex-col bg-[#0d1113]">
      
      <header class="h-20 border-b border-[#1a2326] flex items-center justify-between px-8 bg-[#0d1113]">
        <div class="flex items-center gap-4">
          <h1 class="text-xl font-bold tracking-wide">CHAT EXE101</h1>
        </div>
        <div class="text-xs font-medium text-[#62777a]">
          This is the selected conversation
        </div>
      </header>

      <section class="flex-1 overflow-y-auto px-8 py-6 space-y-6 flex flex-col">
        <div 
          v-for="msg in messages" 
          :key="msg.id"
          :class="[
            msg.role === 'user' ? 'self-end items-end' : 'self-start items-start',
            'flex flex-col max-w-[70%]'
          ]"
        >
          <div 
            v-if="msg.role === 'user'"
            class="bg-[#37736b] text-white px-5 py-3 rounded-[20px] rounded-tr-none text-sm leading-relaxed whitespace-pre-line shadow-md w-fit"
          >
            <p v-for="(line, idx) in msg.text" :key="idx">{{ line }}</p>
          </div>

          <div 
            v-else
            class="text-[#e1e7e8] px-2 py-1 text-sm leading-relaxed whitespace-pre-line w-full"
          >
            <p v-for="(line, idx) in msg.text" :key="idx" class="mb-0.5">{{ line }}</p>
          </div>
        </div>
      </section>

      <footer class="p-8 bg-[#0d1113]">
        <div class="max-w-4xl mx-auto border border-[#2d3a3e] bg-[#13191c] rounded-2xl p-4 flex flex-col gap-3 shadow-lg">
          
          <textarea 
            v-model="inputMessage"
            placeholder="This is a random text to appear in the UI for design"
            class="w-full bg-transparent border-none outline-none resize-none text-sm text-white placeholder-[#4e5f61] h-12 focus:ring-0"
          ></textarea>
          
          <div class="flex items-center justify-between border-t border-[#1e282b] pt-3">
            <button class="flex items-center gap-2 px-3 py-1.5 rounded-full border border-[#2d3a3e] bg-[#1a2326] text-xs font-medium text-[#a6b8ba] hover:bg-[#233034] transition-colors cursor-pointer">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 text-[#4ea399]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
              </svg>
              <span>{{ currentModel }}</span>
            </button>
            
            <button class="w-8 h-8 rounded-full bg-[#1a2326] border border-[#2d3a3e] flex items-center justify-center text-[#a6b8ba] hover:text-[#4ea399] hover:bg-[#233034] transition-colors cursor-pointer">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
              </svg>
            </button>
          </div>

        </div>
      </footer>

    </main>
  </div>
</template>

<style scoped>
/* Bạn có thể thêm animation transition cho các bong bóng chat tại đây nếu muốn */
</style>