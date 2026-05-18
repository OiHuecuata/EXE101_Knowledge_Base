<script setup lang="ts">
import { ref } from 'vue'

interface Message {
  role: 'user' | 'assistant'
  content: string
  sources?: string[]
}

const inputMessage = ref('')
const isStreaming = ref(false)
const messages = ref<Message[]>([
  { 
    role: 'assistant', 
    content: 'Xin chào! Tôi là AI trợ lý tri thức của môn EXE101. Bạn cần hỏi gì về tài liệu hoặc các Case Study?', 
    sources: [] 
  }
])

const sendMessage = () => {
  if (!inputMessage.value.trim() || isStreaming.value) return
  
  // User Message
  messages.value.push({ role: 'user', content: inputMessage.value.trim() })
  const userQuery = inputMessage.value.trim()
  inputMessage.value = ''
  
  // Loading
  isStreaming.value = true
  
  setTimeout(() => {
    messages.value.push({
      role: 'assistant',
      content: `Đây là câu trả lời mẫu cho câu hỏi: "${userQuery}". Hệ thống RAG kết nối Backend Go và Python đang được chuẩn bị kết nối ở bước tiếp theo.`,
      sources: ['W3_Case study 2_Khiem.md', 'W1_case study 1_Khiem.md']
    })
    isStreaming.value = false
  }, 1000)
}
</script>

<template>
  <div class="flex h-screen w-screen bg-slate-900 text-slate-100 font-sans overflow-hidden">
    
    <!-- Sidebar -->
    <aside class="w-72 bg-slate-950 border-r border-slate-800 flex flex-col p-4 space-y-4 select-none">
      <div class="font-bold text-lg border-b border-slate-800 pb-3 flex items-center gap-2">
        <span>📂 EXE101 KB</span>
      </div>
      <div class="flex-1 overflow-y-auto space-y-4">
        <div>
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">Tài liệu đã Ingest</p>
          <div class="space-y-1.5">
            <div class="text-xs p-2 bg-slate-900 hover:bg-slate-850 rounded border border-slate-800 truncate cursor-pointer transition">📄 W3_Case study 2_Khiem.md</div>
            <div class="text-xs p-2 bg-slate-900 hover:bg-slate-850 rounded border border-slate-800 truncate cursor-pointer transition">📄 W1_case study 1_Khiem.md</div>
          </div>
        </div>
      </div>
    </aside>

    <!-- Chat Frame -->
    <main class="flex-1 flex flex-col h-full bg-slate-900 relative">
      
      <!-- Context Layout -->
      <div class="flex-1 overflow-y-auto p-6 space-y-6 max-w-4xl w-full mx-auto">
        <div v-for="(msg, index) in messages" :key="index"
             :class="['p-4 rounded-xl max-w-[85%] flex flex-col space-y-2 transition-all', msg.role === 'user' ? 'bg-blue-600 ml-auto text-white' : 'bg-slate-800 mr-auto border border-slate-750']">
          
          <p class="text-sm leading-relaxed whitespace-pre-line">{{ msg.content }}</p>
          
          <!-- Source -->
          <div v-if="msg.sources && msg.sources.length > 0" class="mt-2 pt-2 border-t border-slate-700 text-xs text-slate-400">
            <span class="font-semibold block mb-1 text-slate-300">🔍 Nguồn tham khảo:</span>
            <div class="flex flex-wrap gap-1.5">
              <span v-for="(src, sIdx) in msg.sources" :key="sIdx" class="bg-slate-950 px-2 py-0.5 rounded border border-slate-700 text-[11px]">
                {{ src }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Input Zone -->
      <div class="p-4 border-t border-slate-800 bg-slate-900/50 backdrop-blur">
        <div class="max-w-4xl mx-auto flex gap-2">
          <input v-model="inputMessage" @keyup.enter="sendMessage" :disabled="isStreaming"
                 type="text" placeholder="Hỏi tôi về các case study môn EXE101..." 
                 class="flex-1 bg-slate-950 border border-slate-800 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500 disabled:opacity-50 transition" />
          <button @click="sendMessage" :disabled="isStreaming"
                  class="bg-blue-600 hover:bg-blue-500 text-white px-5 py-3 rounded-xl text-sm font-medium transition disabled:opacity-50 flex items-center justify-center">
            <span v-if="isStreaming" class="animate-pulse">Đang xử lý...</span>
            <span v-else>Gửi</span>
          </button>
        </div>
      </div>

    </main>
  </div>
</template>