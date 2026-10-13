# Graph Report - app-listener  (2026-10-13)

## Corpus Check
- 369 files · ~1,563,655 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4713 nodes · 19210 edges · 200 communities (172 shown, 28 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1371 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `786e980b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- BinaryEntry
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- update.go
- planUpdaters
- NetGuardEvent
- Candidate
- install.go
- Vault
- IntegrationSuite
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- buildOneGuard
- guard.bpf.c
- NetworkMonitor
- find_test.go
- Ledger
- hashFileShared
- Monitor
- Serve
- runNetworkMonitor
- IntegrationSuite
- fileedit_model_test.go
- IntegrationSuite
- monitorUnitTest
- tempGrantSetup
- controlServer
- bpfstats/main.go
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- IntegrationSuite
- liveSession
- IntegrationSuite
- runUpdate
- highlight_test.go
- .syncSetsLocked
- MountVouch
- tempRow
- inode_dev
- daemon_multicall_test.go
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Cmd
- guardModel
- codeedit.go
- sync.Mutex
- engine
- GuardEvent
- auditGroupCoverage
- filevault_test.go
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- Vault
- find.go
- Config
- ResolveLibraryClosure
- fcntl
- copyRegularAt
- runNetworkGuard
- watchSet
- btrfsLayoutFrom
- time.Time
- github.com/spf13/cobra.Command
- filevault.go
- process_vm_readv.c
- launchKey
- writeWithin
- Read
- update_test.go
- GuardInodeKey
- runDaemon
- render.go
- CleanupStalePins
- IntegrationSuite
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LibraryClosure
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- safeio.go
- .judgeProcess
- IntegrationSuite
- Well-known bypass classes
- golang.org/x/sys/unix.Stat_t
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- CandidateDir
- Build & sign release assets (reusable workflow)
- runMonitor
- harnessSuite
- net_tester/main.go
- NewEventFanout
- app-listener Terminal Demo Recording
- IntegrationSuite
- time.Duration
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- catalogRefresher
- NewDumpHook
- WebSocket /ws client
- netModel
- string
- install.sh
- mc_tag.h
- changelog.go
- make test-integration (rootful Docker suite)
- EventType
- supersede_probe.c
- Run
- trace-app-libs.sh
- IntegrationSuite
- KernelDev
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TemporaryGrant
- RawTarget
- BinaryRule
- preload_wait.c
- Guard
- fscrypt.go
- Vault
- BinaryStat
- uninstall_system.go
- StripPinnedTempAllows
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- Guard
- inodeChain
- New
- vouchedDevs
- Installation
- os.File
- graphify-refresh.sh
- TempBinary
- bufio.Reader
- TrustGuard
- testKeyPair
- bottomBar
- mustSubkey
- integrationtests/guard_test.go
- SuperblockDev
- openBunTmp
- releaseAssets
- lockAndDeprovision
- Limitations
- Troubleshooting
- procRec
- downloadFile
- Daemon configuration
- edit-protected
- absPath
- binaryVetter
- Compatibility
- How it works
- runForward
- markReady
- bpfLsmListed
- unlockUnderGuard
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- parseStableVersion
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- sanityCheckBinary
- fakeDaemon
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS
- hashCache

## God Nodes (most connected - your core abstractions)
1. `Guard` - 95 edges
2. `Config` - 95 edges
3. `GuardInodeKey` - 94 edges
4. `Load()` - 89 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 79 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 66 edges
9. `fileEditModel` - 60 edges
10. `User` - 56 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `Well-known bypass classes` --references--> `btrfs_search POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `Well-known bypass classes` --references--> `io_uring POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Shared bubbletea TUI layout across monitor and guard modes** — media_demo_monitor_tui, media_demo_guard_whitelist_tui, media_demo_guard_blacklist_tui, media_demo_tui_resource_footer [INFERRED 0.85]
- **Reserved-name trust defenses (trust object #3)** — integrationtests_exploits_readme_glob_plant, integrationtests_exploits_readme_lib_probe, readme_bun_tmpdir_redirect, claude_catalog_refresh [INFERRED 0.85]
- **Release build, sign and verify pipeline** — _github_workflows_release_prerelease_workflow, _github_workflows_promote_release_promote_release_workflow, _github_workflows__build_release_build_release_workflow, _github_workflows__build_release_ed25519_signing, readme_install_sh, readme_update [INFERRED 0.95]
- **Verifier budget measurement and gating** — _github_workflows_verifier_verifier_gate, claude_bpfstats, readme_daemon_check, claude_verifier_complexity_budget, _github_workflows__build_release_ebpf_regeneration [INFERRED 0.95]

## Communities (200 total, 28 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.11
Nodes (26): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+18 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (127): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+119 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.05
Nodes (14): IntegrationSuite, IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix() (+6 more)

### Community 3 - "BinaryEntry"
Cohesion: 0.08
Nodes (20): deferredBinary, BinaryEntry, checkBtrfsKeys(), BinariesSummary(), canonicalBinaryPath(), canonicalPaths(), confinedEntry(), eventMask() (+12 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (79): updateCatalogConfig(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+71 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.14
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.09
Nodes (56): trustManager, askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), pickUsers(), addKeysToAgent() (+48 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (21): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), runTool() (+13 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.16
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 11 - "planUpdaters"
Cohesion: 0.12
Nodes (19): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+11 more)

### Community 12 - "NetGuardEvent"
Cohesion: 0.11
Nodes (17): NetGuardEvent, NetworkGuardRepository, NewGuardUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo() (+9 more)

### Community 13 - "Candidate"
Cohesion: 0.09
Nodes (34): appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), pathCovered(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig() (+26 more)

### Community 14 - "install.go"
Cohesion: 0.05
Nodes (74): reloadDaemonForPasswordChange(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), writeEditPasswordHash(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled() (+66 more)

### Community 15 - "Vault"
Cohesion: 0.16
Nodes (15): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+7 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.12
Nodes (3): IntegrationSuite, guardDeltaEvents(), le64HexKey()

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "GenerateConf"
Cohesion: 0.05
Nodes (57): logRefreshChange(), TestLogRefreshChange(), TestDiffMergePreservesExistingAndParses(), askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestGeneratedLibraryBlockRoundTrips(), validateConfigText() (+49 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.26
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (47): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+39 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (9): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "auth.go"
Cohesion: 0.11
Nodes (32): auditAfterEdit(), Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode() (+24 more)

### Community 24 - "buildOneGuard"
Cohesion: 0.13
Nodes (23): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), selfGuards, selfProtectSpec, AdmissionCheck (+15 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (53): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with(), event_is_read_class() (+45 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.07
Nodes (17): NetBpfEvent, NetEvent, NetEventType, Cstr(), FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType() (+9 more)

### Community 27 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 28 - "Ledger"
Cohesion: 0.10
Nodes (17): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), EnsurePlaceholders(), Ledger (+9 more)

### Community 29 - "hashFileShared"
Cohesion: 0.20
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 30 - "Monitor"
Cohesion: 0.11
Nodes (6): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), Monitor, pathCache, pathCacheEntry

### Community 31 - "Serve"
Cohesion: 0.08
Nodes (28): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+20 more)

### Community 32 - "runNetworkMonitor"
Cohesion: 0.28
Nodes (7): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NewNetModel()

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.10
Nodes (41): fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), openRootT(), selectNamed() (+33 more)

### Community 35 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 37 - "tempGrantSetup"
Cohesion: 0.33
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 38 - "controlServer"
Cohesion: 0.22
Nodes (5): controlServer, readClientLines(), streamEvents(), editControlSession, grantRequest

### Community 39 - "bpfstats/main.go"
Cohesion: 0.16
Nodes (20): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc(), TestKernelHasFunc() (+12 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 43 - "liveSession"
Cohesion: 0.24
Nodes (7): parseErr(), displayEntries(), holdForward(), pingInterval(), liveSession, RunFileEditorSession(), RunSession()

### Community 44 - "IntegrationSuite"
Cohesion: 0.17
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "runUpdate"
Cohesion: 0.16
Nodes (19): changelogText(), showChangelog(), TestChangelogText(), applyUpdate(), confirmUpdate(), downloadAndVerify(), filterChannel(), isTerminal() (+11 more)

### Community 46 - "highlight_test.go"
Cohesion: 0.11
Nodes (18): diffModel, Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+10 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "MountVouch"
Cohesion: 0.07
Nodes (19): GlobChild, statFunc, deleteResKeys(), syncU32Map(), compileGlobNames(), globText(), TrustGuard, resolveBits() (+11 more)

### Community 50 - "tempRow"
Cohesion: 0.29
Nodes (4): tempGrant, tempMaskOp, tempRow, liveTrust()

### Community 51 - "inode_dev"
Cohesion: 0.14
Nodes (19): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), code_suspect(), current_is_inspector() (+11 more)

### Community 52 - "daemon_multicall_test.go"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.14
Nodes (37): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+29 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.11
Nodes (9): clamp(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), assertViewFits(), TestViewsFitReportedTerminalSize(), noticeModel, servedModel, serveTestModel (+1 more)

### Community 55 - "guardModel"
Cohesion: 0.17
Nodes (9): formatGuardEventLine(), formatDecision(), formatGuardType(), listenForGuardEvents(), tickGuardStats(), syncViewport(), guardEventLine, guardEventMsg (+1 more)

### Community 56 - "codeedit.go"
Cohesion: 0.25
Nodes (16): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), Navigate(), onlySpaces() (+8 more)

### Community 57 - "sync.Mutex"
Cohesion: 0.17
Nodes (6): controlManager, controlManager, newControlManager(), startControlManager(), startControlServer(), configEditor

### Community 58 - "engine"
Cohesion: 0.09
Nodes (9): TestTieRank(), engine, pathWithin(), TestPathWithinMatchesPrefixRule(), tieRank(), deleteInoKeys(), applyInodeCapacity(), intersectAllows() (+1 more)

### Community 59 - "GuardEvent"
Cohesion: 0.07
Nodes (23): holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), bpfGuardEvent, TestWithChmodDropWriteNeedsReadOnly() (+15 more)

### Community 60 - "auditGroupCoverage"
Cohesion: 0.20
Nodes (12): auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+4 more)

### Community 61 - "filevault_test.go"
Cohesion: 0.29
Nodes (18): looksLikeFileVaultRecord(), inode(), TestEnsureRecoverySidecarPlaceholder(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestStageRecoverySealsContentAndStaysInPlace(), TestUnlockFileInPlaceRecoversFromInterruptedTransform() (+10 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.18
Nodes (17): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+9 more)

### Community 64 - "errno"
Cohesion: 0.28
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.15
Nodes (25): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+17 more)

### Community 66 - "Vault"
Cohesion: 0.24
Nodes (6): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata()

### Community 67 - "find.go"
Cohesion: 0.19
Nodes (8): findButton, span, findAll(), lowerRunes(), newFindInput(), overlay(), runesEqual(), TestFindAllSmartCase()

### Community 68 - "Config"
Cohesion: 0.17
Nodes (18): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), trustManager, inGuardedTree(), keys() (+10 more)

### Community 69 - "ResolveLibraryClosure"
Cohesion: 0.20
Nodes (10): fileExists(), LdSoPreloadPaths(), ResolveLibraryClosure(), skipWithoutHostRoot(), TestMountVouchesOwnershipJudgesHostNamespace(), TestResolveLibraryClosure(), TestResolveLibraryClosureForeignArch(), TestResolveLibraryClosureScript() (+2 more)

### Community 70 - "fcntl"
Cohesion: 0.16
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "runNetworkGuard"
Cohesion: 0.15
Nodes (14): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), eventsetSummary(), Mode (+6 more)

### Community 73 - "watchSet"
Cohesion: 0.14
Nodes (14): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+6 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "time.Time"
Cohesion: 0.16
Nodes (15): newGateLogLimiter(), gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter() (+7 more)

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (18): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+10 more)

### Community 78 - "filevault.go"
Cohesion: 0.24
Nodes (13): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace() (+5 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 81 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 82 - "Read"
Cohesion: 0.15
Nodes (14): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+6 more)

### Community 83 - "update_test.go"
Cohesion: 0.12
Nodes (7): fetchReleases(), readInstalledVersion(), rsaPublicKeyPEM(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestParsePublicKey(), TestReadInstalledVersion()

### Community 84 - "GuardInodeKey"
Cohesion: 0.06
Nodes (28): refusedReplacements, appendInodeRow(), batchInodeRows(), inodeRowsOf(), engine, Inspector, SetInspectors(), engine (+20 more)

### Community 85 - "runDaemon"
Cohesion: 0.10
Nodes (23): CheckBPFLSM(), CheckEBPF(), configureDaemonLogging(), loadDaemonConfig(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath(), runBPFCheck() (+15 more)

### Community 86 - "render.go"
Cohesion: 0.15
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LibraryClosure"
Cohesion: 0.16
Nodes (12): ldConfParser, defaultLibDirs(), elfInterp(), isScript(), ldConfInclude(), ldSoConfDirs(), LibraryClosure(), resolveSoname() (+4 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.11
Nodes (27): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless() (+19 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.14
Nodes (7): FirstOnBtrfs(), OnBtrfs(), SuperblockDevPath(), TestSuperblockDevMatchesStatOffBtrfs(), NetGuard, mcNameKey(), statInodeKey()

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "safeio.go"
Cohesion: 0.06
Nodes (43): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+35 more)

### Community 100 - ".judgeProcess"
Cohesion: 0.20
Nodes (11): reasonText(), codeVouched(), NetGuard, isTraced(), loaderEnv(), roFor(), TestCodeVouched(), TestLoaderEnv() (+3 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "CandidateDir"
Cohesion: 0.06
Nodes (32): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder (+24 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runMonitor"
Cohesion: 0.11
Nodes (15): MakeDisplayPaths(), ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule(), runHeadless() (+7 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "NewEventFanout"
Cohesion: 0.22
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "time.Duration"
Cohesion: 0.11
Nodes (11): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), debouncer, New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod() (+3 more)

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "catalogRefresher"
Cohesion: 0.19
Nodes (6): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), catalogRefresher

### Community 121 - "NewDumpHook"
Cohesion: 0.14
Nodes (15): applyVerbosity(), VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys(), sortStrings(), TestBridgeStdLog() (+7 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "netModel"
Cohesion: 0.12
Nodes (13): sanitizeTerminalText(), sanitizeTerminalTexts(), clipToWidth(), formatResourceBar(), formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents() (+5 more)

### Community 124 - "string"
Cohesion: 0.16
Nodes (6): main(), report(), copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 127 - "changelog.go"
Cohesion: 0.25
Nodes (6): newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes(), changelogModel

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "EventType"
Cohesion: 0.08
Nodes (12): BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), planTempBlock(), Run() (+4 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "KernelDev"
Cohesion: 0.21
Nodes (11): fdStat, physicalParent(), statFD(), binaryStatOf(), KernelDev(), LstatInode(), BtrfsDevices(), btrfsSysfsDevices() (+3 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TemporaryGrant"
Cohesion: 0.15
Nodes (8): TemporaryGrant, TempRule, applyTempGrants(), TemporaryAccess, planTempGrants(), revokeTempGrants(), tempJournalRows(), fakeTempGrant

### Community 143 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 144 - "BinaryRule"
Cohesion: 0.29
Nodes (8): systemPatterns(), systemRule(), grantNewBinaries(), Line(), Lines(), TestLinesKeyOnLinkAndIncludePending(), BinaryRule, ParseEventType()

### Community 146 - "Guard"
Cohesion: 0.09
Nodes (5): rootHandle, openRoot(), SharedPinDegraded(), Guard, NewGuard()

### Community 147 - "fscrypt.go"
Cohesion: 0.05
Nodes (60): confirmMasterKeyOverwrite(), runGenKey(), TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs() (+52 more)

### Community 149 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 150 - "uninstall_system.go"
Cohesion: 0.08
Nodes (35): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), removeKeyAndEmptyDir(), removeMasterKey(), removeSSHAgentEnv(), removeSSHAgentUnit() (+27 more)

### Community 151 - "StripPinnedTempAllows"
Cohesion: 0.38
Nodes (3): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.17
Nodes (29): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+21 more)

### Community 155 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 156 - "New"
Cohesion: 0.26
Nodes (11): editorModel, moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation() (+3 more)

### Community 157 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.04
Nodes (64): multicallAdmissible(), ownerPaths(), checkFreshDir(), openDirNoFollow(), renewDir(), setConfinedHomes(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted() (+56 more)

### Community 161 - "TempBinary"
Cohesion: 0.31
Nodes (8): TempBinary, closeTempBinaries(), daemonUseCase, logTempGrant(), newExeTap(), resolveTempBinaries(), sanitized(), watchTempBinaries()

### Community 162 - "bufio.Reader"
Cohesion: 0.18
Nodes (12): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 163 - "TrustGuard"
Cohesion: 0.09
Nodes (10): trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied() (+2 more)

### Community 164 - "testKeyPair"
Cohesion: 0.27
Nodes (8): parseChecksum(), parsePublicKey(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestVerifyRelease(), verifyRelease(), verifyReleaseWithKey()

### Community 165 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 166 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 167 - "integrationtests/guard_test.go"
Cohesion: 0.27
Nodes (6): attrOpCase, eventsForPath(), guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 168 - "SuperblockDev"
Cohesion: 0.22
Nodes (7): BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), SuperblockDev(), TestMountinfoDev(), TestMountinfoHasFstype(), progressReader

### Community 169 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 170 - "releaseAssets"
Cohesion: 0.32
Nodes (8): assetsFor(), checkUpdatePreconditions(), resolveAssets(), TestAssetsFor(), TestResolveAssetsPicksArch(), TestValidateChannel(), validateChannel(), releaseAssets

### Community 171 - "lockAndDeprovision"
Cohesion: 0.33
Nodes (6): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "procRec"
Cohesion: 0.33
Nodes (5): judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), procRec

### Community 175 - "downloadFile"
Cohesion: 0.47
Nodes (4): downloadFile(), downloadReleaseFiles(), TestDownloadFileHTTPError(), TestDownloadFileMode0700()

### Community 176 - "Daemon configuration"
Cohesion: 0.22
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "absPath"
Cohesion: 0.09
Nodes (8): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase, absPath(), TestMain(), helperChild(), TestMain()

### Community 179 - "binaryVetter"
Cohesion: 0.13
Nodes (18): newBinaryVetter(), openBinaryVetter(), backingDeviceUnion(), buildConcurrency(), buildGuards(), makeReloadHandler(), reloadOnce(), reloadSlotsNeeded() (+10 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "runForward"
Cohesion: 0.50
Nodes (5): chooseResources(), eventsLabel(), forwardBinaries(), normalizedEvents(), runForward()

### Community 183 - "markReady"
Cohesion: 0.15
Nodes (13): markReady(), noteReloadedSettings(), notifySystemdReady(), inodeBatchFn, fakeBatches(), TestBatchInodeRows(), engine, ConfigureInodeCapacity() (+5 more)

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 185 - "unlockUnderGuard"
Cohesion: 0.50
Nodes (5): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), SectionScan

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "parseStableVersion"
Cohesion: 0.40
Nodes (5): compareStableVersions(), parseStableVersion(), parseVerPart(), TestParseStableVersion(), stableVersion

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 194 - "sanityCheckBinary"
Cohesion: 0.40
Nodes (4): sanityCheckBinary(), TestSanityCheckBinaryRefusesOtherMachine(), runsNatively(), TestRunsNatively()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.21
Nodes (7): TestFilterExistingWhitelistSymlinkEscape(), applyNewFileMeta(), defaultNewFileMode(), listEntries(), TestDefaultNewFileMode(), writeFileInPlace(), vaultFS

### Community 201 - "hashCache"
Cohesion: 0.18
Nodes (8): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache, Ledger

## Knowledge Gaps
- **147 isolated node(s):** `One-line installer`, `Build from source`, `SSH agent`, `Bun-based apps (opencode)`, `fscrypt prerequisites` (+142 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 488 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **28 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Why does `fileEditModel` connect `fileEditModel` to `fileedit_model_test.go`, `vaultFS`, `github.com/charmbracelet/bubbles/textarea.Model`, `update.go`, `highlight_test.go`, `github.com/charmbracelet/bubbletea.Cmd`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Why does `fakeGuardRepo` connect `fakeGuardRepo` to `daemon/daemon.go`, `TempBinary`, `TemporaryGrant`, `usecase/daemon_test.go`, `sync.Mutex`, `GuardEvent`, `os.File`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **What connects `One-line installer`, `Build from source`, `SSH agent` to the rest of the system?**
  _147 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10869831922463502 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02289803220035778 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.05165289256198347 - nodes in this community are weakly interconnected._