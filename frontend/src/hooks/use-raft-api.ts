"use client"

import { useState, useEffect, useCallback } from "react"
import type { RaftNode, ApiResponse } from "@/lib/raft-types"

const API_BASE_URL = "http://localhost:8080/api"

export function useRaftApi() {
  const [nodes, setNodes] = useState<RaftNode[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchNodes = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/nodes`)
      if (!response.ok) {
        throw new Error("Failed to fetch nodes")
      }
      const result: ApiResponse<RaftNode[]> = await response.json();
      
      if (result.error) {
        throw new Error(result.error.message);
      }
      
      setNodes(result.data || [])
      setLoading(false)
    } catch (err) {
      setError("Error fetching nodes")
      setLoading(false)
    }
  }, [])

  const startElection = useCallback(
    async (nodeId: number) => {
      try {
        const response = await fetch(`${API_BASE_URL}/nodes/${nodeId}/start-election`, {
          method: "POST",
        })
        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || "启动选举失败");
        }
        await fetchNodes() // Refresh nodes after starting election
      } catch (err) {
        console.error('选举操作失败:', err);
        setError(err instanceof Error ? err.message : "发生未知错误");
      }
    },
    [fetchNodes],
  )

  useEffect(() => {
    fetchNodes()
  }, [fetchNodes])

  return { nodes, loading, error, fetchNodes, startElection }
}
