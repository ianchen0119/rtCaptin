# RT-Captin

RT-Captin is the worker pool implementation based on real-time mechanism.
It's inspired by the Priority Inheritance Protocol (PIP) & Preemption concept in OS.

## Why RT-Captin?
- Prevent deadlock due to synchronization issues.
- Suitable for network service, which contains the multiple message types
- Suitable for concurrent program

## Key concepts
- The lower priority's job will be preempted if it is preemptable
  - This job will be terminated by receiving the early-break event from a specific channel
- Otherwise, the lower priority's job will inherit the priority from another job, which tries to use the shared resource it has.