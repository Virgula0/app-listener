# Graph Report - app-listener  (2026-10-13)

## Corpus Check
- 369 files · ~1,563,655 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4709 nodes · 19245 edges · 193 communities (164 shown, 29 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1371 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `eacaef5f`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- orphans_test.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- go_pkg_path_filepath
- planUpdaters
- TestMonitorUseCaseLifecycle
- Candidate
- install.go
- wipe
- IntegrationSuite
- guard_trust.bpf.c
- SetSectionWhitelist
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- Hash
- buildOneGuard
- guard.bpf.c
- NetworkMonitor
- find_test.go
- binaryVetter
- BinaryEntry
- Monitor
- Serve
- runNetworkGuard
- IntegrationSuite
- fileedit_model_test.go
- IntegrationSuite
- dataRowRenderer
- tempGrantSetup
- sync.Mutex
- bpfstats/main.go
- github.com/charmbracelet/bubbles/textarea.Model
- .SetGlobReservations
- absPath
- liveSession
- IntegrationSuite
- update.go
- highlight_test.go
- .syncSetsLocked
- TrustGuard
- tempRow
- inode_dev
- main_test.go
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Cmd
- daemonModel
- codeedit.go
- controlManager
- engine
- GuardEvent
- audit.go
- filevault_test.go
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- Vault
- find.go
- buildTrustedSet
- ResolveLibraryClosure
- fcntl
- copyRegularAt
- netGuardModel
- time.Time
- btrfsLayoutFrom
- GenerateConf
- github.com/spf13/cobra.Command
- IntegrationSuite
- process_vm_readv.c
- launchKey
- writeWithin
- Read
- newTestVetter
- GuardInodeKey
- runLockdown
- render.go
- .write
- fakeVault
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LibraryClosure
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- unit.go
- KernelDev
- IntegrationSuite
- Well-known bypass classes
- NetworkMonitorUseCase
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- CandidateDir
- Build & sign release assets (reusable workflow)
- runMonitor
- fileedit_symlink_test.go
- net_tester/main.go
- NewEventFanout
- app-listener Terminal Demo Recording
- IntegrationSuite
- Pulse
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- time.Duration
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
- ParseEventsFlag
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- DiscoverForUsers
- trustManager
- preload_wait.c
- Guard
- systemd.go
- Vault
- uninstall_system.go
- StripPinnedTempAllows
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- IntegrationSuite
- inodeChain
- highlight.go
- vouchedDevs
- Installation
- os.File
- graphify-refresh.sh
- bufio.Reader
- github.com/cilium/ebpf/link.Link
- IntegrationSuite
- is_event_type_allowed
- parseGuardEvent
- Limitations
- Troubleshooting
- Daemon configuration
- edit-protected
- IntegrationSuite
- Vault
- Compatibility
- How it works
- runForward
- loadDaemonConfig
- bpfLsmListed
- launch_rule_at
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- newWatchPattern
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- fakeDaemon
- findTopLevelWindow
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS
- hashCache
- countInodeRows

## God Nodes (most connected - your core abstractions)
1. `Config` - 95 edges
2. `Guard` - 95 edges
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

## Communities (193 total, 29 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (20): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+12 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (114): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+106 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.08
Nodes (11): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+3 more)

### Community 3 - "orphans_test.go"
Cohesion: 0.21
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (78): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+70 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.06
Nodes (75): trustManager, refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup() (+67 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (21): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), runTool() (+13 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.25
Nodes (25): buildGlobReservations(), bitOf(), discordLibConfig(), hasChild(), mkdirs(), steamLibDirConfig(), TestBuildGlobReservations_BunRootsListed(), TestBuildGlobReservations_BunTmpdir() (+17 more)

### Community 11 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 12 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.21
Nodes (13): NewGuardUseCase(), NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError() (+5 more)

### Community 13 - "Candidate"
Cohesion: 0.09
Nodes (30): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+22 more)

### Community 14 - "install.go"
Cohesion: 0.08
Nodes (38): applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), deploy() (+30 more)

### Community 15 - "wipe"
Cohesion: 0.13
Nodes (17): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn() (+9 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "SetSectionWhitelist"
Cohesion: 0.07
Nodes (38): logRefreshChange(), TestLogRefreshChange(), askEncryption(), LibraryBlock, RefreshOptions, findSectionEnd(), findSectionStart(), firstContentLine() (+30 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (43): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+35 more)

### Community 22 - "Resource"
Cohesion: 0.09
Nodes (11): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles(), uniqueEncryptionRoots() (+3 more)

### Community 23 - "Hash"
Cohesion: 0.13
Nodes (28): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+20 more)

### Community 24 - "buildOneGuard"
Cohesion: 0.09
Nodes (36): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), resolveGuardConfig(), runGuard(), runGuardTUI() (+28 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (51): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), evict_inode_from_guard(), get_dentry_from_path() (+43 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.11
Nodes (10): NetBpfEvent, NetEvent, FormatAddr(), Ntohs(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey() (+2 more)

### Community 27 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 28 - "binaryVetter"
Cohesion: 0.08
Nodes (22): newBinaryVetter(), openBinaryVetter(), candidates(), confirmable(), confirmAll(), describe(), hashNow(), run() (+14 more)

### Community 29 - "BinaryEntry"
Cohesion: 0.07
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), BinariesSummary(), canonicalPaths() (+21 more)

### Community 30 - "Monitor"
Cohesion: 0.06
Nodes (8): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), Monitor, monitorUnitTest, pathCache, pathCacheEntry

### Community 31 - "Serve"
Cohesion: 0.08
Nodes (28): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+20 more)

### Community 32 - "runNetworkGuard"
Cohesion: 0.12
Nodes (22): CheckEBPF(), ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard() (+14 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.12
Nodes (32): AtomicWriteAt(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth() (+24 more)

### Community 35 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 37 - "tempGrantSetup"
Cohesion: 0.21
Nodes (10): tempJournalRows(), journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure() (+2 more)

### Community 38 - "sync.Mutex"
Cohesion: 0.17
Nodes (8): controlManager, controlServer, readClientLines(), startControlServer(), streamEvents(), configEditor, editControlSession, grantRequest

### Community 39 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 40 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.26
Nodes (3): column(), Finder, overlay()

### Community 41 - ".SetGlobReservations"
Cohesion: 0.16
Nodes (10): GlobChild, statFunc, compileGlobNames(), globText(), TrustGuard, resolveBits(), resolveChildren(), statKey() (+2 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "IntegrationSuite"
Cohesion: 0.18
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "update.go"
Cohesion: 0.06
Nodes (47): showChangelog(), applyUpdate(), assetsFor(), checkUpdatePreconditions(), compareStableVersions(), confirmUpdate(), downloadAndVerify(), downloadFile() (+39 more)

### Community 46 - "highlight_test.go"
Cohesion: 0.18
Nodes (24): moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines() (+16 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "TrustGuard"
Cohesion: 0.07
Nodes (14): trustRows, deleteResKeys(), setSupersedeTrust(), syncU32Map(), cStr(), syncMap(), TrustGuard, logTrustDenied() (+6 more)

### Community 50 - "tempRow"
Cohesion: 0.17
Nodes (7): exeRow, tempGrant, tempMaskOp, tempRow, Guard, liveTrust(), planTempBlock()

### Community 51 - "inode_dev"
Cohesion: 0.12
Nodes (23): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), code_suspect(), current_is_inspector() (+15 more)

### Community 52 - "main_test.go"
Cohesion: 0.05
Nodes (20): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+12 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.14
Nodes (39): sameConfig(), watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector() (+31 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.11
Nodes (7): clamp(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 55 - "daemonModel"
Cohesion: 0.26
Nodes (6): listenForDaemonEvents(), tickDaemonStats(), syncViewport(), daemonEventLine, daemonEventMsg, daemonModel

### Community 56 - "codeedit.go"
Cohesion: 0.25
Nodes (16): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), Navigate(), onlySpaces() (+8 more)

### Community 57 - "controlManager"
Cohesion: 0.29
Nodes (3): controlManager, newControlManager(), startControlManager()

### Community 58 - "engine"
Cohesion: 0.08
Nodes (10): TestTieRank(), engine, pathWithin(), TestPathWithinMatchesPrefixRule(), tieRank(), deleteInoKeys(), guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent() (+2 more)

### Community 59 - "GuardEvent"
Cohesion: 0.09
Nodes (14): holdHeadless(), HeadlessLine(), runGuardHeadless(), GuardEvent, logBacklogDenial(), formatGuardEventLine(), formatDecision(), formatGuardType() (+6 more)

### Community 60 - "audit.go"
Cohesion: 0.23
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 61 - "filevault_test.go"
Cohesion: 0.11
Nodes (41): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, looksLikeFileVaultRecord(), newFileVaultAEAD() (+33 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+10 more)

### Community 63 - "fakeGuardRepo"
Cohesion: 0.09
Nodes (3): vetDecision, SupersededBinary, fakeGuardRepo

### Community 64 - "errno"
Cohesion: 0.28
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.17
Nodes (23): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+15 more)

### Community 66 - "Vault"
Cohesion: 0.20
Nodes (9): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow(), TestOpenVaultRefusesSymlink() (+1 more)

### Community 67 - "find.go"
Cohesion: 0.27
Nodes (7): findButton, span, findAll(), lowerRunes(), newFindInput(), runesEqual(), TestFindAllSmartCase()

### Community 68 - "buildTrustedSet"
Cohesion: 0.16
Nodes (16): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+8 more)

### Community 69 - "ResolveLibraryClosure"
Cohesion: 0.18
Nodes (11): hostView(), mountVouchesOwnership(), ResolveLibraryClosure(), skipWithoutHostRoot(), TestHostView(), TestMountVouchesOwnershipJudgesHostNamespace(), TestResolveLibraryClosure(), TestResolveLibraryClosureForeignArch() (+3 more)

### Community 70 - "fcntl"
Cohesion: 0.16
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "netGuardModel"
Cohesion: 0.13
Nodes (7): NetGuardEvent, NetworkGuardRepository, protoString(), NetworkGuardUseCase, netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 73 - "time.Time"
Cohesion: 0.15
Nodes (11): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID, sweepMemo (+3 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "GenerateConf"
Cohesion: 0.10
Nodes (28): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), sectionsFromCandidates(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+20 more)

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (18): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+10 more)

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
Cohesion: 0.10
Nodes (19): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+11 more)

### Community 83 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 84 - "GuardInodeKey"
Cohesion: 0.05
Nodes (32): refusedReplacements, appendInodeRow(), batchInodeRows(), inodeRowsOf(), engine, Inspector, TrustGuard, SetInspectors() (+24 more)

### Community 85 - "runLockdown"
Cohesion: 0.18
Nodes (12): CheckBPFLSM(), configureDaemonLogging(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runLockdown(), runOneShotMode(), dirIsRootOwnedSafe() (+4 more)

### Community 86 - "render.go"
Cohesion: 0.13
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - ".write"
Cohesion: 0.15
Nodes (12): rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParsePublicKey(), TestVerifyRelease(), CleanupStalePins(), genOfPinFile(), TestCleanupStalePins() (+4 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LibraryClosure"
Cohesion: 0.12
Nodes (16): ldConfParser, defaultLibDirs(), elfInterp(), isScript(), ldConfInclude(), ldSoConfDirs(), LibraryClosure(), rootOwnedInode() (+8 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.11
Nodes (32): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+24 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.13
Nodes (9): NetEventType, NetEventTypes(), ParseNetEventType(), eventsetSummary(), eventTypeKey(), NetGuard, mcNameKey(), NewNetGuard() (+1 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "unit.go"
Cohesion: 0.18
Nodes (14): TestDaemonUnitWithMetadataOutput(), DaemonUnitWithMetadataOutput(), DesiredDaemonUnit(), ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary() (+6 more)

### Community 100 - "KernelDev"
Cohesion: 0.05
Nodes (48): fdStat, BinaryStat, placedWalk, physicalParent(), statFD(), binaryStatOf(), binaryStatOf(), KernelDev() (+40 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

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
Cohesion: 0.07
Nodes (25): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), confSafeMatch() (+17 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runMonitor"
Cohesion: 0.11
Nodes (19): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+11 more)

### Community 110 - "fileedit_symlink_test.go"
Cohesion: 0.17
Nodes (12): harnessSuite, mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent() (+4 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "NewEventFanout"
Cohesion: 0.22
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "Pulse"
Cohesion: 0.23
Nodes (7): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "time.Duration"
Cohesion: 0.13
Nodes (8): drain(), newCatalogRefresher(), readEvents(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), catalogRefresher, debouncer

### Community 121 - "NewDumpHook"
Cohesion: 0.14
Nodes (15): applyVerbosity(), VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys(), sortStrings(), TestBridgeStdLog() (+7 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "netModel"
Cohesion: 0.18
Nodes (10): runTUI(), formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), NewNetModel(), tickNetStats(), netEventLine (+2 more)

### Community 124 - "string"
Cohesion: 0.16
Nodes (6): main(), report(), copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "changelog.go"
Cohesion: 0.22
Nodes (8): changelogText(), newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes(), changelogModel

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "EventType"
Cohesion: 0.10
Nodes (11): BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), Run(), Cstr() (+3 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "Run"
Cohesion: 0.13
Nodes (8): columnState, dataRow, guiModel, headerWidget, newColumnState(), newDataRow(), newHeaderWidget(), Run()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.14
Nodes (15): forwardRequest, TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), daemonUseCase, TemporaryAccess (+7 more)

### Community 143 - "DiscoverForUsers"
Cohesion: 0.33
Nodes (7): Discover(), DiscoverForUsers(), DiscoverSystem(), TestDiscoverForUsersDeduplicates(), TestDiscoverOnlyExisting(), TestDiscoverPerUserHomes(), TestDiscoverSystemOnlyAbsolute()

### Community 144 - "trustManager"
Cohesion: 0.18
Nodes (11): resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), trustManager, startTrustGuard(), trustStartupError(), grantNewBinaries(), inspectorAdmitter, Line() (+3 more)

### Community 146 - "Guard"
Cohesion: 0.07
Nodes (13): deferredBinary, rootHandle, openRoot(), SharedPinDegraded(), canonicalBinaryPath(), eventMask(), Guard, patternMatches() (+5 more)

### Community 147 - "systemd.go"
Cohesion: 0.05
Nodes (62): reloadDaemonForPasswordChange(), confirmRunPrereq(), resolveFilesystemPrereqs(), revertSystemFiles(), checkKeyLen(), classifySetupError(), classifySupportError(), GenerateMasterKey() (+54 more)

### Community 150 - "uninstall_system.go"
Cohesion: 0.05
Nodes (60): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), cleanOrphanedFscrypt(), deletePostBackups(), restoreBackups(), askDecrypt(), cleanOrphanedMetadata() (+52 more)

### Community 151 - "StripPinnedTempAllows"
Cohesion: 0.38
Nodes (3): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.24
Nodes (21): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+13 more)

### Community 155 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 156 - "highlight.go"
Cohesion: 0.14
Nodes (6): diffModel, editorModel, Highlighter, lexerFor(), shebangLexer(), tokenClassOf()

### Community 157 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.04
Nodes (69): multicallAdmissible(), ownerPaths(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun() (+61 more)

### Community 162 - "bufio.Reader"
Cohesion: 0.15
Nodes (12): swapFile(), controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines() (+4 more)

### Community 163 - "github.com/cilium/ebpf/link.Link"
Cohesion: 0.26
Nodes (4): trustHook, foreignPinnedLink(), NetGuard, hook

### Community 166 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 169 - "parseGuardEvent"
Cohesion: 0.22
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.22
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 179 - "Vault"
Cohesion: 0.10
Nodes (30): backingDeviceUnion(), buildConcurrency(), buildGuards(), lockOneRoot(), lockRootRecovering(), makeReloadHandler(), newPinGeneration(), prepareDaemonStart() (+22 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "runForward"
Cohesion: 0.18
Nodes (13): chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents(), pingInterval() (+5 more)

### Community 183 - "loadDaemonConfig"
Cohesion: 0.15
Nodes (12): loadDaemonConfig(), markReady(), noteReloadedSettings(), notifySystemdReady(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), engine (+4 more)

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 185 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "newWatchPattern"
Cohesion: 0.17
Nodes (18): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+10 more)

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

### Community 195 - "fakeDaemon"
Cohesion: 0.15
Nodes (4): systemPatterns(), systemRule(), fakeDaemon, ParseEventType()

### Community 196 - "findTopLevelWindow"
Cohesion: 0.38
Nodes (4): findTopLevelWindow(), floatWindow(), internAtom(), setFloatHint()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.31
Nodes (4): applyNewFileMeta(), listEntries(), dirEntry, vaultFS

### Community 201 - "hashCache"
Cohesion: 0.18
Nodes (8): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache, Ledger

### Community 205 - "countInodeRows"
Cohesion: 0.67
Nodes (4): inodeBatchFn, fakeBatches(), countInodeRows(), TestCountInodeRows()

## Knowledge Gaps
- **147 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+142 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 484 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **29 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.051) - this node is a cross-community bridge._
- **Why does `engine` connect `engine` to `daemon/daemon.go`, `github.com/cilium/ebpf/link.Link`, `sync.Mutex`, `bpfstats/main.go`, `Vault`, `GuardInodeKey`, `NetGuard`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `fakeDaemon`, `sync.Mutex`, `time.Time`, `Vault`, `Pulse`, `GuardInodeKey`, `buildOneGuard`, `controlManager`, `GuardEvent`, `BinaryEntry`, `os.File`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _147 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11797184638513052 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.025346301208370173 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08302469135802469 - nodes in this community are weakly interconnected._