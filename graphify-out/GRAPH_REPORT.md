# Graph Report - app-listener  (2026-10-12)

## Corpus Check
- 359 files · ~1,558,243 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4641 nodes · 18683 edges · 182 communities (157 shown, 25 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1344 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9c5a3e1a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- supersedeMaps
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- setpassword.go
- planUpdaters
- TestMonitorUseCaseLifecycle
- Candidate
- install.go
- fscrypt.go
- .startGuardStd
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- update.go
- buildOneGuard
- guard.bpf.c
- NetworkMonitor
- find_test.go
- Ledger
- netGuardModel
- Monitor
- Serve
- GuardInodeKey
- IntegrationSuite
- fileedit_model_test.go
- IntegrationSuite
- app.go
- protected/daemon_test.go
- controlServer
- bpfstats/main.go
- Finder
- TempBinary
- IntegrationSuite
- liveSession
- IntegrationSuite
- update_test.go
- highlight_test.go
- .syncSetsLocked
- TrustGuard
- tempRow
- inode_dev
- IntegrationSuite
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Cmd
- absPath
- codeedit.go
- controlManager
- ResolveLibraryClosure
- downloadFile
- auditGroupCoverage
- runNetworkGuard
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- uninstall_system.go
- CandidateDir
- filevault.go
- Pulse
- fcntl
- copyRegularAt
- trustManager
- is_event_type_allowed
- btrfsLayoutFrom
- monitorUnitTest
- github.com/spf13/cobra.Command
- IntegrationSuite
- bottomBar
- process_vm_readv.c
- golang.org/x/sys/unix.Stat_t
- BinaryEntry
- Read
- releaseAssets
- MountVouch
- sync.Mutex
- github.com/charmbracelet/bubbles/textarea.Model
- watchSet
- newerThanStable
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LibraryClosure
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- dialLiveSession
- KernelDev
- IntegrationSuite
- Well-known bypass classes
- ResolvePinBase
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- tempGrant
- Build & sign release assets (reusable workflow)
- runMonitor
- fileedit_symlink_test.go
- net_tester/main.go
- NewEventFanout
- app-listener Terminal Demo Recording
- IntegrationSuite
- journalTo
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- time.Duration
- runDaemon
- WebSocket /ws client
- netModel
- string
- install.sh
- mc_tag.h
- changelog.go
- make test-integration (rootful Docker suite)
- FileEvent
- supersede_probe.c
- bufio.Reader
- trace-app-libs.sh
- IntegrationSuite
- sanityCheckBinary
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- .GrantTemporaryAccess
- launch_rule_at
- SanitizeText
- preload_wait.c
- Guard
- writeWithin
- launchKey
- offline.go
- daemon: permanent protection with encryption at rest
- newTestVetter
- New
- Installation
- os.File
- graphify-refresh.sh
- GuardEvent
- loadDaemonConfig
- github.com/cilium/ebpf/link.Link
- parseGuardEvent
- Limitations
- Troubleshooting
- inodeChain
- integrationtests/guard_test.go
- Daemon configuration
- edit-protected
- NewDumpHook
- Compatibility
- How it works
- time.Time
- bpfLsmListed
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS

## God Nodes (most connected - your core abstractions)
1. `Config` - 95 edges
2. `Guard` - 95 edges
3. `GuardInodeKey` - 92 edges
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

## Communities (182 total, 25 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.11
Nodes (25): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+17 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (111): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+103 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.05
Nodes (14): IntegrationSuite, IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix() (+6 more)

### Community 3 - "supersedeMaps"
Cohesion: 0.09
Nodes (19): checkBtrfsKeys(), ensureBtrfsLayout(), dropSupersededLocked(), Guard, PruneSuperseded(), supersededMark(), supersedeMaps(), supersedePruneLoop() (+11 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (83): TestSetSectionWhitelistPreservesGroupStructure(), updateCatalogConfig(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines() (+75 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.05
Nodes (85): trustManager, askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), restoreDaemonAfterDiffAbort() (+77 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (18): guardUnitTest, ReplacementCheck, guardModeKey(), addRecorder(), copySelf(), runTool(), walkInodes(), walkLiveEntries() (+10 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.16
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "setpassword.go"
Cohesion: 0.07
Nodes (48): ParseEventsFlag(), Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode() (+40 more)

### Community 11 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 12 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.38
Nodes (7): NewGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestNetworkGuardUseCaseLifecycle(), fakeNetworkGuardRepo

### Community 13 - "Candidate"
Cohesion: 0.08
Nodes (39): appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), pathCovered(), runDiffCatalog(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates() (+31 more)

### Community 14 - "install.go"
Cohesion: 0.07
Nodes (51): applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), encryptDirectories(), ensureInstalledBinary(), ensureMasterKey(), installBinaryOnly() (+43 more)

### Community 15 - "fscrypt.go"
Cohesion: 0.06
Nodes (41): confirmMasterKeyOverwrite(), runGenKey(), askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty(), deprovisionKind (+33 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (52): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+44 more)

### Community 18 - "GenerateConf"
Cohesion: 0.05
Nodes (58): logRefreshChange(), TestLogRefreshChange(), TestDiffMergePreservesExistingAndParses(), askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestGeneratedLibraryBlockRoundTrips() (+50 more)

### Community 19 - "fileEditModel"
Cohesion: 0.10
Nodes (9): clipLabel(), humanSize(), clamp(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.16
Nodes (46): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+38 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (14): ownerPaths(), fakeDaemon, Resource, GuardRepository, Vault, concurrencyLimit(), encryptionRootOf(), DaemonUseCase (+6 more)

### Community 23 - "update.go"
Cohesion: 0.21
Nodes (16): applyUpdate(), confirmUpdate(), downloadAndVerify(), filterChannel(), isTerminal(), latestTime(), newerThanInstalled(), pickLatestRelease() (+8 more)

### Community 24 - "buildOneGuard"
Cohesion: 0.09
Nodes (36): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), resolveGuardConfig(), runGuard(), runGuardTUI() (+28 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (53): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with(), event_is_read_class() (+45 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.07
Nodes (18): NetBpfEvent, NetEvent, NetEventType, FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType(), NewNetworkMonitor() (+10 more)

### Community 27 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 28 - "Ledger"
Cohesion: 0.09
Nodes (18): grantNewBinaries(), candidates(), confirmable(), confirmAll(), describe(), run(), EnsurePlaceholders(), Ledger (+10 more)

### Community 29 - "netGuardModel"
Cohesion: 0.22
Nodes (5): NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 30 - "Monitor"
Cohesion: 0.12
Nodes (5): evalSymlinksOrEmpty(), NewMonitor(), Monitor, pathCache, pathCacheEntry

### Community 31 - "Serve"
Cohesion: 0.08
Nodes (28): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+20 more)

### Community 32 - "GuardInodeKey"
Cohesion: 0.05
Nodes (27): refusedReplacements, TestNearestRootClaims(), TestTieRank(), appendInodeRow(), batchInodeRows(), engine, inodeRowsOf(), nearestRootClaims() (+19 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.15
Nodes (27): fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd() (+19 more)

### Community 35 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 36 - "app.go"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 37 - "protected/daemon_test.go"
Cohesion: 0.12
Nodes (26): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), removeKeyAndEmptyDir(), removeMasterKey(), TestRemoveKeyAndEmptyDir(), Delete() (+18 more)

### Community 38 - "controlServer"
Cohesion: 0.18
Nodes (8): controlManager, controlServer, readClientLines(), startControlServer(), streamEvents(), configEditor, editControlSession, grantRequest

### Community 39 - "bpfstats/main.go"
Cohesion: 0.14
Nodes (23): CheckBPFLSM(), runBPFCheck(), guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity() (+15 more)

### Community 40 - "Finder"
Cohesion: 0.15
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 41 - "TempBinary"
Cohesion: 0.21
Nodes (7): TempBinary, TemporaryGrant, TempRule, planTempBlock(), newExeTap(), planTempGrants(), watchTempBinaries()

### Community 44 - "IntegrationSuite"
Cohesion: 0.17
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "update_test.go"
Cohesion: 0.13
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 46 - "highlight_test.go"
Cohesion: 0.11
Nodes (18): editorModel, Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+10 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "TrustGuard"
Cohesion: 0.09
Nodes (12): statFunc, trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), setSupersedeTrust(), cStr(), resolveBits() (+4 more)

### Community 50 - "tempRow"
Cohesion: 0.30
Nodes (4): exeRow, tempMaskOp, tempRow, Guard

### Community 51 - "inode_dev"
Cohesion: 0.14
Nodes (19): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), code_suspect(), current_is_inspector() (+11 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.06
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.14
Nodes (40): sameConfig(), watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector() (+32 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.13
Nodes (9): listenForGuardEvents(), tickGuardStats(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), guardEventMsg, noticeModel, servedModel, serveTestModel (+1 more)

### Community 55 - "absPath"
Cohesion: 0.08
Nodes (7): IntegrationSuite, IntegrationSuite, launchCase, absPath(), TestMain(), helperChild(), TestMain()

### Community 56 - "codeedit.go"
Cohesion: 0.25
Nodes (17): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), Navigate() (+9 more)

### Community 57 - "controlManager"
Cohesion: 0.29
Nodes (3): controlManager, newControlManager(), startControlManager()

### Community 58 - "ResolveLibraryClosure"
Cohesion: 0.14
Nodes (14): hostView(), mountVouchesOwnership(), ResolveLibraryClosure(), rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), skipWithoutHostRoot(), TestHostView() (+6 more)

### Community 59 - "downloadFile"
Cohesion: 0.22
Nodes (8): downloadFile(), downloadReleaseFiles(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleases(), TestFetchReleasesHTTPError(), progressReader

### Community 60 - "auditGroupCoverage"
Cohesion: 0.16
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 61 - "runNetworkGuard"
Cohesion: 0.11
Nodes (22): CheckEBPF(), ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard() (+14 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.15
Nodes (20): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+12 more)

### Community 64 - "errno"
Cohesion: 0.28
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.10
Nodes (34): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), runLockdown(), runOneShotMode(), ensurePinStateFilePlaceholder(), ensurePlaceholder() (+26 more)

### Community 66 - "uninstall_system.go"
Cohesion: 0.31
Nodes (4): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), hint

### Community 67 - "CandidateDir"
Cohesion: 0.06
Nodes (31): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, GlobChild, compileGlobNames(), globText() (+23 more)

### Community 68 - "filevault.go"
Cohesion: 0.09
Nodes (45): harnessSuite, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, looksLikeFileVaultRecord() (+37 more)

### Community 69 - "Pulse"
Cohesion: 0.19
Nodes (8): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub, EventFanout

### Community 70 - "fcntl"
Cohesion: 0.16
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "trustManager"
Cohesion: 0.13
Nodes (20): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), trustManager (+12 more)

### Community 73 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "monitorUnitTest"
Cohesion: 0.10
Nodes (3): newPathCache(), dirTarget(), monitorUnitTest

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (19): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 78 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.14
Nodes (16): BinaryStat, placedWalk, binaryStatOf(), binaryStatOf(), StatBinary(), StatBinaryFile(), confineBelow(), inUserHome() (+8 more)

### Community 81 - "BinaryEntry"
Cohesion: 0.10
Nodes (20): hashNow(), binaryVerifyState, sharedHashEntry, BinaryEntry, BinariesSummary(), canonicalPaths(), hashFileShared(), hashPathShared() (+12 more)

### Community 82 - "Read"
Cohesion: 0.10
Nodes (19): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+11 more)

### Community 83 - "releaseAssets"
Cohesion: 0.32
Nodes (8): assetsFor(), checkUpdatePreconditions(), resolveAssets(), TestAssetsFor(), TestResolveAssetsPicksArch(), TestValidateChannel(), validateChannel(), releaseAssets

### Community 84 - "MountVouch"
Cohesion: 0.10
Nodes (12): deleteResKeys(), syncU32Map(), loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows(), syncMap(), deleteKeys(), MountVouch (+4 more)

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.16
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - "watchSet"
Cohesion: 0.19
Nodes (10): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchDir, watchPattern (+2 more)

### Community 88 - "newerThanStable"
Cohesion: 0.33
Nodes (7): compareStableVersions(), newerThanStable(), parseStableVersion(), parseVerPart(), TestNewerThanStable(), TestParseStableVersion(), stableVersion

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LibraryClosure"
Cohesion: 0.14
Nodes (14): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+6 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.09
Nodes (35): UIDResolver, NewUIDResolver(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI(), writeEvent() (+27 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.14
Nodes (6): reasonText(), eventTypeKey(), NetGuard, mcNameKey(), statInodeKey(), TestReasonText()

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "dialLiveSession"
Cohesion: 0.16
Nodes (15): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+7 more)

### Community 100 - "KernelDev"
Cohesion: 0.08
Nodes (30): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+22 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runMonitor"
Cohesion: 0.11
Nodes (19): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+11 more)

### Community 110 - "fileedit_symlink_test.go"
Cohesion: 0.35
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "NewEventFanout"
Cohesion: 0.25
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout[T]

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "journalTo"
Cohesion: 0.53
Nodes (6): journalTo(), tempBinary(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure()

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "time.Duration"
Cohesion: 0.13
Nodes (8): drain(), newCatalogRefresher(), readEvents(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), catalogRefresher, debouncer

### Community 121 - "runDaemon"
Cohesion: 0.09
Nodes (29): awaitStartupOrSignal(), backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), drainEphemeral(), makeReloadHandler(), markReady() (+21 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "netModel"
Cohesion: 0.16
Nodes (11): runTUI(), syncViewport(), formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), NewNetModel(), tickNetStats() (+3 more)

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
Cohesion: 0.21
Nodes (9): changelogText(), newChangelogModel(), showChangelog(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "FileEvent"
Cohesion: 0.09
Nodes (9): BpfEvent, eventUnitTest, FileEvent, Run(), Cstr(), DecodeBpfEvent(), newFakeMonitorRepo(), editorHarness (+1 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "bufio.Reader"
Cohesion: 0.14
Nodes (13): swapFile(), controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines() (+5 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "sanityCheckBinary"
Cohesion: 0.40
Nodes (4): sanityCheckBinary(), TestSanityCheckBinaryRefusesOtherMachine(), runsNatively(), TestRunsNatively()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - ".GrantTemporaryAccess"
Cohesion: 0.20
Nodes (9): applyTempGrants(), closeTempBinaries(), daemonUseCase, TemporaryAccess, logTempGrant(), resolveTempBinaries(), revokeTempGrants(), sanitized() (+1 more)

### Community 143 - "launch_rule_at"
Cohesion: 0.50
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 144 - "SanitizeText"
Cohesion: 0.09
Nodes (17): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), systemPatterns(), systemRule(), binaryVetter, vetDecision, EventType (+9 more)

### Community 146 - "Guard"
Cohesion: 0.06
Nodes (16): deferredBinary, rootHandle, openRoot(), SharedPinDegraded(), canonicalBinaryPath(), confinedEntry(), eventMask(), Guard (+8 more)

### Community 147 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 148 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 150 - "offline.go"
Cohesion: 0.06
Nodes (50): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), TestCollectFilesystemPrereqsSkipsRegularFiles(), classifySetupError() (+42 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "newTestVetter"
Cohesion: 0.18
Nodes (22): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+14 more)

### Community 156 - "New"
Cohesion: 0.26
Nodes (11): diffModel, moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation() (+3 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.05
Nodes (63): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+55 more)

### Community 165 - "GuardEvent"
Cohesion: 0.09
Nodes (12): holdHeadless(), HeadlessLine(), runGuardHeadless(), commMatchesGuardedBinary(), GuardEvent, logBacklogDenial(), formatGuardEventLine(), formatDecision() (+4 more)

### Community 167 - "loadDaemonConfig"
Cohesion: 0.29
Nodes (7): loadDaemonConfig(), noteReloadedSettings(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 169 - "github.com/cilium/ebpf/link.Link"
Cohesion: 0.19
Nodes (8): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile(), NetGuard, hook

### Community 170 - "parseGuardEvent"
Cohesion: 0.22
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 175 - "integrationtests/guard_test.go"
Cohesion: 0.21
Nodes (7): attrOpCase, eventsForPath(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 176 - "Daemon configuration"
Cohesion: 0.22
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 179 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 183 - "time.Time"
Cohesion: 0.14
Nodes (13): newWatchSet(), fakeInotify, inodeBatchFn, sweepID, sweepMemo, fakeBatches(), TestBatchInodeRows(), engine (+5 more)

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 186 - "Binary trust"
Cohesion: 0.33
Nodes (6): Binary ledger, Binary replacement, Binary trust, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

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

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.11
Nodes (17): addErrKind, classifyAddErr(), walkEntries(), TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), applyNewFileMeta(), defaultNewFileMode() (+9 more)

## Knowledge Gaps
- **147 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+142 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 484 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **25 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `integrationtests/guard_test.go`, `.startGuardStd`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.052) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `GuardInodeKey`, `Pulse`, `GuardEvent`, `SanitizeText`, `BinaryEntry`, `sync.Mutex`, `Resource`, `time.Time`, `buildOneGuard`, `runDaemon`, `controlManager`, `os.File`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Why does `Monitor` connect `Monitor` to `daemon/daemon.go`, `github.com/cilium/ebpf/link.Link`, `NetworkMonitor`, `sync.Mutex`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _147 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11169210351692103 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.026238161776121722 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.052556615866685655 - nodes in this community are weakly interconnected._