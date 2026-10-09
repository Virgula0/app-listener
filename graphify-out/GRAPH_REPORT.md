# Graph Report - app-listener  (2026-10-09)

## Corpus Check
- 345 files · ~1,549,271 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4535 nodes · 17924 edges · 198 communities (168 shown, 30 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1294 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f6db6701`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- filevault.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- mkdirs
- auth.go
- uninstall_system.go
- codeedit.go
- BinaryEntry
- systemd.go
- Vault
- IntegrationSuite
- guard_trust.bpf.c
- confgen.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- DaemonUseCase
- RawTarget
- Monitor
- guard.bpf.c
- dataRowRenderer
- github.com/spf13/cobra.Command
- binaryVetter
- buildOneGuard
- highlight_test.go
- serve.go
- engine
- time.Duration
- fileedit_model_test.go
- DaemonEvent
- watchGroup
- filevault_test.go
- GuardEvent
- guardSpec
- github.com/charmbracelet/bubbles/textarea.Model
- NetEventType
- absPath
- liveSession
- IntegrationSuite
- fscrypt/orphans.go
- monitorUnitTest
- .syncSetsLocked
- IntegrationSuite
- tempRow
- __always_inline
- IntegrationSuite
- TrustGuard
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- os.File
- controlServer
- install.go
- EventType
- vaultFS
- runNetworkGuard
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- netModel
- highlight.go
- planUpdaters
- catalogRefresher
- fcntl
- copytree.go
- bufio.Reader
- app.go
- btrfsLayoutFrom
- Ledger
- RefreshSection
- FileEvent
- process_vm_readv.c
- Read
- uninstall.go
- update_test.go
- CandidateDir
- NetworkMonitor
- safeio.go
- render.go
- time.Time
- TempBinary
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- KernelDev
- auditGroupCoverage
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- serveWebSocket
- sanitizeTerminalText
- IntegrationSuite
- Well-known bypass classes
- Vault
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- github.com/cilium/ebpf/link.Link
- Build & sign release assets (reusable workflow)
- newTestVetter
- Candidate
- net_tester/main.go
- TestMonitorUseCaseLifecycle
- app-listener Terminal Demo Recording
- is_event_type_allowed
- Run
- lockAndDeprovision
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- mustSubkey
- buildTrustedSet
- WebSocket /ws client
- NewDumpHook
- launch_probe.c
- install.sh
- mc_tag.h
- Config
- make test-integration (rootful Docker suite)
- NewEventFanout
- supersede_probe.c
- stripTempJournal
- trace-app-libs.sh
- IntegrationSuite
- sync.Mutex
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- bottomBar
- find_test.go
- runMonitor
- Guard
- GuardInodeKey
- update.go
- configedit_test.go
- validateFlags
- writeWithin
- daemon: permanent protection with encryption at rest
- dialLiveSession
- IntegrationSuite
- tempGrantSetup
- launchKey
- inodeChain
- Installation
- searchDirs
- graphify-refresh.sh
- downloadFile
- runLockdown
- changelog.go
- DiscoverForUsers
- fileedit_symlink_test.go
- launch_rule_at
- find.go
- Vault
- bpfstats/main.go
- watchPattern
- Limitations
- Troubleshooting
- headerWidget
- fakeDaemon
- Daemon configuration
- edit-protected
- IntegrationSuite
- loadDaemonConfig
- Compatibility
- How it works
- io.Reader
- StripPinnedTempAllows
- bpfLsmListed
- .pathInfo
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- SuperblockDev
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- harnessSuite
- TestMain
- .runEditorHarness
- monitor: see what touches a path

## God Nodes (most connected - your core abstractions)
1. `Config` - 93 edges
2. `Guard` - 93 edges
3. `GuardInodeKey` - 86 edges
4. `Load()` - 85 edges
5. `absPath()` - 78 edges
6. `IntegrationSuite` - 75 edges
7. `IntegrationSuite` - 66 edges
8. `writeConfig()` - 65 edges
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

## Communities (198 total, 30 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.09
Nodes (39): ownerPaths(), backingDeviceUnion(), libBinaries(), resourceBinaries(), resourceLinks(), Execute(), tempJournalRow, classCacheKey (+31 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (114): plantBun(), TestOpenBunTmp_KeepsVettedDir(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait() (+106 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "filevault.go"
Cohesion: 0.23
Nodes (15): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace() (+7 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (79): TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines() (+71 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.09
Nodes (57): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), deploy(), preflightDeployedBinary(), cleanOrphanedFscrypt() (+49 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (21): addErrKind, guardUnitTest, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), runTool() (+13 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "auth.go"
Cohesion: 0.08
Nodes (38): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+30 more)

### Community 11 - "uninstall_system.go"
Cohesion: 0.32
Nodes (7): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), GeneralToolsOf(), namedGeneralTool(), ResolvesToGeneralTool(), hint

### Community 12 - "codeedit.go"
Cohesion: 0.25
Nodes (16): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), Navigate(), onlySpaces() (+8 more)

### Community 13 - "BinaryEntry"
Cohesion: 0.10
Nodes (16): deferredBinary, VettedInode, BinaryEntry, BinariesSummary(), canonicalBinaryPath(), canonicalPaths(), ConfinedEntry(), confinedEntry() (+8 more)

### Community 14 - "systemd.go"
Cohesion: 0.05
Nodes (67): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), reloadDaemonForPasswordChange(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), revertSystemFiles() (+59 more)

### Community 15 - "Vault"
Cohesion: 0.16
Nodes (15): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+7 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "confgen.go"
Cohesion: 0.09
Nodes (41): LibraryBlock, Section, findSectionEnd(), findSectionStart(), firstContentLine(), GenerateConf(), GenerateSections(), isLibDirectiveText() (+33 more)

### Community 19 - "fileEditModel"
Cohesion: 0.10
Nodes (9): clipLabel(), humanSize(), clamp(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.16
Nodes (45): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+37 more)

### Community 22 - "DaemonUseCase"
Cohesion: 0.18
Nodes (5): concurrencyLimit(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles(), unlockRoots()

### Community 23 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (46): add_inode_to_guard(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), evict_inode_from_guard(), fill_path(), get_dentry_from_path(), get_inode_from_path() (+38 more)

### Community 27 - "github.com/spf13/cobra.Command"
Cohesion: 0.10
Nodes (28): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+20 more)

### Community 28 - "binaryVetter"
Cohesion: 0.12
Nodes (20): newBinaryVetter(), openBinaryVetter(), buildConcurrency(), buildGuards(), makeReloadHandler(), reloadOnce(), reloadSlotsNeeded(), startCatalogRefresh() (+12 more)

### Community 29 - "buildOneGuard"
Cohesion: 0.14
Nodes (22): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), selfGuards, AdmissionCheck, BinaryRule, libDirWriters(), TestDispatchWithoutOwnAllowedEvents() (+14 more)

### Community 30 - "highlight_test.go"
Cohesion: 0.18
Nodes (24): moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines() (+16 more)

### Community 31 - "serve.go"
Cohesion: 0.11
Nodes (22): basicAuth(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop(), relayServeSignals() (+14 more)

### Community 32 - "engine"
Cohesion: 0.08
Nodes (9): TestTieRank(), engine, pathWithin(), tieRank(), deleteInoKeys(), guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), intersectAllows() (+1 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.15
Nodes (27): fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd() (+19 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.09
Nodes (32): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+24 more)

### Community 36 - "watchGroup"
Cohesion: 0.11
Nodes (35): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+27 more)

### Community 37 - "filevault_test.go"
Cohesion: 0.30
Nodes (17): looksLikeFileVaultRecord(), inode(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsFileVaultRecordShape(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceRecoversFromInterruptedTransform(), TestUnlockFileInPlaceWrongKeyFails() (+9 more)

### Community 38 - "GuardEvent"
Cohesion: 0.06
Nodes (26): selfProtectSpecs(), displayEntries(), holdForward(), holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless() (+18 more)

### Community 39 - "guardSpec"
Cohesion: 0.25
Nodes (10): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), InertProgram(), KernelHasFunc(), TestKernelHasFunc(), trustSpec() (+2 more)

### Community 40 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.26
Nodes (3): column(), Finder, overlay()

### Community 41 - "NetEventType"
Cohesion: 0.10
Nodes (9): NetBpfEvent, NetEvent, NetEventType, FormatAddr(), Ntohs(), ParseNetEventType(), NetworkMonitorRepository, NetworkMonitorUseCase (+1 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "IntegrationSuite"
Cohesion: 0.18
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "fscrypt/orphans.go"
Cohesion: 0.19
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 46 - "monitorUnitTest"
Cohesion: 0.09
Nodes (4): NewMonitor(), newPathCache(), dirTarget(), monitorUnitTest

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 50 - "tempRow"
Cohesion: 0.17
Nodes (7): exeRow, tempGrant, tempMaskOp, tempRow, Guard, liveTrust(), planTempBlock()

### Community 51 - "__always_inline"
Cohesion: 0.11
Nodes (27): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), btrfs_copy_gate(), check_and_emit_args(), code_suspect() (+19 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.06
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 53 - "TrustGuard"
Cohesion: 0.07
Nodes (15): trustHook, trustRows, deleteResKeys(), setSupersedeTrust(), syncU32Map(), cStr(), syncMap(), TrustGuard (+7 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.10
Nodes (12): listenForGuardEvents(), tickGuardStats(), syncViewport(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), guardEventLine, guardEventMsg, guardModel (+4 more)

### Community 56 - "os.File"
Cohesion: 0.05
Nodes (60): multicallAdmissible(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched() (+52 more)

### Community 57 - "controlServer"
Cohesion: 0.22
Nodes (5): controlServer, readClientLines(), streamEvents(), editControlSession, grantRequest

### Community 58 - "install.go"
Cohesion: 0.09
Nodes (32): TestDiffMergePreservesExistingAndParses(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), encryptDirectories(), ensureInstalledBinary(), ensureMasterKey(), installBinaryOnly() (+24 more)

### Community 59 - "EventType"
Cohesion: 0.09
Nodes (13): systemPatterns(), systemRule(), BpfEvent, EventType, eventUnitTest, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+5 more)

### Community 60 - "vaultFS"
Cohesion: 0.14
Nodes (13): TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), applyNewFileMeta(), defaultNewFileMode(), listEntries(), openRootT(), TestDefaultNewFileMode() (+5 more)

### Community 61 - "runNetworkGuard"
Cohesion: 0.10
Nodes (17): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), Mode (+9 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.19
Nodes (17): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+9 more)

### Community 64 - "errno"
Cohesion: 0.21
Nodes (5): put(), wait_then_read(), dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.29
Nodes (14): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+6 more)

### Community 66 - "netModel"
Cohesion: 0.17
Nodes (11): runNetworkMonitor(), runTUI(), formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), NewNetModel(), tickNetStats() (+3 more)

### Community 67 - "highlight.go"
Cohesion: 0.14
Nodes (6): diffModel, editorModel, Highlighter, lexerFor(), shebangLexer(), tokenClassOf()

### Community 68 - "planUpdaters"
Cohesion: 0.08
Nodes (30): hardLinked(), mayUpdate(), planUpdaters(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning(), TestIsGeneralTool(), TestIsGeneralToolJudgesTheInode(), TestMayUpdate_HardLinkWaivedForUserWritableLibBinary() (+22 more)

### Community 69 - "catalogRefresher"
Cohesion: 0.16
Nodes (8): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), catalogRefresher, Parse(), validateResources()

### Community 70 - "fcntl"
Cohesion: 0.13
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copytree.go"
Cohesion: 0.20
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "bufio.Reader"
Cohesion: 0.18
Nodes (12): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 73 - "app.go"
Cohesion: 0.16
Nodes (9): columnState, dataRow, findTopLevelWindow(), floatWindow(), internAtom(), newColumnState(), newDataRow(), newHeaderWidget() (+1 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "Ledger"
Cohesion: 0.10
Nodes (18): grantNewBinaries(), candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), EnsurePlaceholders() (+10 more)

### Community 76 - "RefreshSection"
Cohesion: 0.15
Nodes (16): logRefreshChange(), TestLogRefreshChange(), RefreshOptions, SectionScan, catalogEntryMatch(), findCatalogRoot(), findCatalogWatchSubPath(), isInsidePath() (+8 more)

### Community 78 - "FileEvent"
Cohesion: 0.13
Nodes (6): runHeadless(), FileEvent, Run(), MonitorRepository, MonitorUseCase, fakeMonitorRepo

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "Read"
Cohesion: 0.15
Nodes (14): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+6 more)

### Community 81 - "uninstall.go"
Cohesion: 0.10
Nodes (32): deletePostBackups(), restoreBackups(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), offerBackupCleanup(), pickDirsToDecrypt() (+24 more)

### Community 82 - "update_test.go"
Cohesion: 0.09
Nodes (17): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled(), TestNewerThanStable() (+9 more)

### Community 83 - "CandidateDir"
Cohesion: 0.05
Nodes (36): trustManager, entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, GlobChild, statFunc (+28 more)

### Community 84 - "NetworkMonitor"
Cohesion: 0.18
Nodes (6): NetEventTypes(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 85 - "safeio.go"
Cohesion: 0.08
Nodes (32): TestDaemonUnitWithMetadataOutput(), updateCatalogConfig(), daemonUnitWithMetadataOutput(), dupFD(), installConfig(), installFile(), installFileAs(), installServices() (+24 more)

### Community 86 - "render.go"
Cohesion: 0.13
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - "time.Time"
Cohesion: 0.19
Nodes (9): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, TestNearestRootClaims(), nearestRootClaims() (+1 more)

### Community 88 - "TempBinary"
Cohesion: 0.13
Nodes (17): TempBinary, TemporaryGrant, TempRule, ResolveTempBinary(), TestResolveTempBinaryClassesRuntime(), applyTempGrants(), closeTempBinaries(), daemonUseCase (+9 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "KernelDev"
Cohesion: 0.06
Nodes (44): fdStat, BinaryStat, placedWalk, physicalParent(), statFD(), binaryStatOf(), isFuseType(), parseMajorMinor() (+36 more)

### Community 92 - "auditGroupCoverage"
Cohesion: 0.16
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.13
Nodes (9): reasonText(), eventsetSummary(), eventTypeKey(), NetGuard, mcNameKey(), modeLabel(), NewNetGuard(), statInodeKey() (+1 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "serveWebSocket"
Cohesion: 0.21
Nodes (7): configureWebSocket(), sameOrigin(), serveWebSocket(), TestSameOrigin(), writeSnapshot(), writeSnapshotLoop(), snapshotHub

### Community 100 - "sanitizeTerminalText"
Cohesion: 0.17
Nodes (8): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), TestSanitizeTerminalText()

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "Vault"
Cohesion: 0.22
Nodes (7): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "github.com/cilium/ebpf/link.Link"
Cohesion: 0.33
Nodes (3): foreignPinnedLink(), NetGuard, hook

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "newTestVetter"
Cohesion: 0.42
Nodes (9): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSwapWithinGenerationLosesRights() (+1 more)

### Community 110 - "Candidate"
Cohesion: 0.22
Nodes (15): appendSectionsAndEdit(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates() (+7 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.19
Nodes (14): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError() (+6 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 116 - "lockAndDeprovision"
Cohesion: 0.33
Nodes (6): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "mustSubkey"
Cohesion: 0.33
Nodes (8): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat(), TestSealOpenFileVaultRoundTrip()

### Community 121 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "NewDumpHook"
Cohesion: 0.14
Nodes (15): applyVerbosity(), VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys(), sortStrings(), TestBridgeStdLog() (+7 more)

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "Config"
Cohesion: 0.10
Nodes (29): newDaemonModel(), runTUI(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), trustManager, startTrustGuard(), trustStartupError(), applyDiffAdditions() (+21 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "NewEventFanout"
Cohesion: 0.22
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "stripTempJournal"
Cohesion: 0.22
Nodes (11): ensurePinStateFilePlaceholder(), ensurePlaceholder(), writeInPlace(), bootstrapSelfProtectPlaceholders(), ensureHashFilePlaceholder(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace() (+3 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "sync.Mutex"
Cohesion: 0.17
Nodes (6): controlManager, controlManager, newControlManager(), startControlManager(), startControlServer(), configEditor

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 143 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 144 - "runMonitor"
Cohesion: 0.20
Nodes (10): CheckEBPF(), ParseEventsFlag(), newPinGeneration(), prepareDaemonStart(), validateForwardRule(), runMonitor(), runTUI(), PrintLogo() (+2 more)

### Community 146 - "Guard"
Cohesion: 0.08
Nodes (12): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, rootHandle, sharedHashEntry, openRoot() (+4 more)

### Community 147 - "GuardInodeKey"
Cohesion: 0.06
Nodes (28): refusedReplacements, engine, Inspector, SetInspectors(), engine, supersedeUnnamed(), appletTag(), buildOf() (+20 more)

### Community 148 - "update.go"
Cohesion: 0.14
Nodes (28): showChangelog(), assetsFor(), checkUpdatePreconditions(), compareStableVersions(), confirmUpdate(), downloadAndVerify(), filterChannel(), isTerminal() (+20 more)

### Community 149 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 150 - "validateFlags"
Cohesion: 0.50
Nodes (4): TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags()

### Community 151 - "writeWithin"
Cohesion: 0.14
Nodes (13): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+5 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "dialLiveSession"
Cohesion: 0.28
Nodes (9): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+1 more)

### Community 155 - "tempGrantSetup"
Cohesion: 0.25
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 156 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 157 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "searchDirs"
Cohesion: 0.22
Nodes (8): ldConfParser, defaultLibDirs(), elfInterp(), ldConfInclude(), ldSoConfDirs(), resolveSoname(), searchDirs(), walkNeeded()

### Community 161 - "downloadFile"
Cohesion: 0.22
Nodes (8): downloadFile(), downloadReleaseFiles(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleases(), TestFetchReleasesHTTPError(), progressReader

### Community 162 - "runLockdown"
Cohesion: 0.17
Nodes (13): CheckBPFLSM(), configureDaemonLogging(), confirmMasterKeyOverwrite(), resolveConfigPath(), runBPFCheck(), runGenKey(), runLockdown(), runOneShotMode() (+5 more)

### Community 163 - "changelog.go"
Cohesion: 0.22
Nodes (8): changelogText(), newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes(), changelogModel

### Community 164 - "DiscoverForUsers"
Cohesion: 0.33
Nodes (7): Discover(), DiscoverForUsers(), DiscoverSystem(), TestDiscoverForUsersDeduplicates(), TestDiscoverOnlyExisting(), TestDiscoverPerUserHomes(), TestDiscoverSystemOnlyAbsolute()

### Community 165 - "fileedit_symlink_test.go"
Cohesion: 0.40
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 166 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 167 - "find.go"
Cohesion: 0.27
Nodes (7): findButton, span, findAll(), lowerRunes(), newFindInput(), runesEqual(), TestFindAllSmartCase()

### Community 169 - "bpfstats/main.go"
Cohesion: 0.40
Nodes (8): chownToSudoUser(), compare(), dumpFullLog(), dumpTrace(), loadOne(), main(), measure(), short()

### Community 170 - "watchPattern"
Cohesion: 0.39
Nodes (6): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchPattern

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 179 - "loadDaemonConfig"
Cohesion: 0.29
Nodes (7): loadDaemonConfig(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), SetUserHomes(), TestResolveConfinedRootOwnedParents(), TestOpenSystemPlacedRefusesUserHome()

### Community 180 - "Compatibility"
Cohesion: 0.29
Nodes (7): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.29
Nodes (7): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`

### Community 182 - "io.Reader"
Cohesion: 0.29
Nodes (6): entryOf(), BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype()

### Community 183 - "StripPinnedTempAllows"
Cohesion: 0.38
Nodes (3): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 185 - ".pathInfo"
Cohesion: 0.48
Nodes (3): evalSymlinksOrEmpty(), pathCache, pathCacheEntry

### Community 186 - "Binary trust"
Cohesion: 0.33
Nodes (6): Binary ledger, Binary replacement, Binary trust, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "SuperblockDev"
Cohesion: 0.33
Nodes (6): checkBtrfsKeys(), FirstOnBtrfs(), OnBtrfs(), SuperblockDev(), SuperblockDevPath(), TestSuperblockDevMatchesStatOffBtrfs()

### Community 190 - "Development"
Cohesion: 0.40
Nodes (5): Demo GIFs, Development, Docker, Makefile targets, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 195 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **140 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+135 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 472 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **30 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `auth.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `DaemonEvent`, `sync.Mutex`, `GuardEvent`, `BinaryEntry`, `fakeDaemon`, `GuardInodeKey`, `time.Time`, `os.File`, `binaryVetter`, `buildOneGuard`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **Why does `shQuote()` connect `IntegrationSuite` to `daemon/daemon.go`, `.runEditorHarness`, `absPath`, `auth.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _140 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09175602444367578 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02549616108938143 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0898838004101162 - nodes in this community are weakly interconnected._