#pragma once

// prctl(PR_SET_MM, ...) sub-options that repoint mm->exe_file: PR_SET_MM_EXE_FILE and PR_SET_MM_MAP
// (which carries an exe_fd). Every guard object keys a caller on mm->exe_file, so without this gate
// a process opens a whitelisted binary and borrows its identity (CRIU's restore primitive, reachable
// unprivileged in a user namespace). A task_prctl hook MUST return -ENOSYS for every other prctl:
// 0 makes prctl() short-circuit, silently no-op'ing legit options (PR_SET_NAME, ...).
#define PR_SET_MM 35
#define PR_SET_MM_EXE_FILE 13
#define PR_SET_MM_MAP 14

static __always_inline int prctl_gate(unsigned long long *ctx)
{
	int option = (int)ctx[0];
	unsigned long arg2 = (unsigned long)ctx[1];
	if (option == PR_SET_MM && (arg2 == PR_SET_MM_MAP || arg2 == PR_SET_MM_EXE_FILE))
		return -EPERM;
	return -ENOSYS;
}
